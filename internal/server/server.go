package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"time"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"
	secretscontroller "github.com/FlutterDizaster/EncryNest/internal/server/controllers/secrets-controller"
	usercontroller "github.com/FlutterDizaster/EncryNest/internal/server/controllers/user-controller"
	jwtresolver "github.com/FlutterDizaster/EncryNest/internal/server/jwt-resolver"
	"github.com/FlutterDizaster/EncryNest/internal/server/repository/postgres"
	secretsrepo "github.com/FlutterDizaster/EncryNest/internal/server/repository/secrets"
	usersrepo "github.com/FlutterDizaster/EncryNest/internal/server/repository/users"
	"github.com/FlutterDizaster/EncryNest/internal/server/services/interceptors"
	secretsservice "github.com/FlutterDizaster/EncryNest/internal/server/services/secrets"
	userservice "github.com/FlutterDizaster/EncryNest/internal/server/services/users"
	"github.com/FlutterDizaster/EncryNest/pkg/keychain"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	TokenTTL         = time.Hour * 24 * 31
	DBConnRetryCount = 3
	DBConnRetryDelay = time.Second
)

type Settings struct {
	Addr                string `desc:"Server address"  env:"SERVER_ADDR"    name:"addr"       short:"a"`
	Port                string `desc:"Server port"     env:"SERVER_PORT"    name:"port"       short:"p" default:"50555"`
	JWTSecret           string `desc:"JWT secret"      env:"JWT_SECRET"     name:"jwt-secret" short:"s"`
	DatabaseURL         string `desc:"Database URL"    env:"DATABASE_URL"   name:"db-url"     short:"d"`
	MigrationsDirectory string `desc:"Migrations dir"  env:"MIGRATIONS_DIR" name:"migrations" short:"m"`
	CertsDirectory      string `desc:"Certs directory" env:"CERTS_DIR"      name:"certs-dir"  short:"e"`
	CertFileName        string `desc:"Cert file name"  env:"CERT_FILE_NAME" name:"cert-file"  short:"c" default:"cert.pem"`
	KeyFileName         string `desc:"Key file name"   env:"KEY_FILE_NAME"  name:"key-file"   short:"k" default:"key.pem"`
}

type Server struct {
	addr           string
	port           string
	jwtSecret      string
	databaseURL    string
	migrationsDir  string
	certsDirectory string
	certFileName   string
	keyFileName    string
	srv            *grpc.Server
}

func NewServer(settings Settings) *Server {
	return &Server{
		addr:           settings.Addr,
		port:           settings.Port,
		jwtSecret:      settings.JWTSecret,
		databaseURL:    settings.DatabaseURL,
		migrationsDir:  settings.MigrationsDirectory,
		certsDirectory: settings.CertsDirectory,
		certFileName:   settings.CertFileName,
		keyFileName:    settings.KeyFileName,
	}
}

func (s *Server) Init(ctx context.Context) error {
	slog.Info("Starting server")

	// Setup JWT resolver
	jwtResolverSettings := jwtresolver.Settings{
		Secret:   s.jwtSecret,
		TokenTTL: TokenTTL,
	}
	jwtResolver := jwtresolver.New(jwtResolverSettings)

	// Setup repositories
	var userRepo usercontroller.UserRepository
	var secretsRepo secretscontroller.SecretsRepository

	if s.databaseURL == "" {
		userRepo = usersrepo.NewInMemoryRepository()
		secretsRepo = secretsrepo.NewInMemoryRepository()
	} else {
		// Setup database
		poolManager, err := postgres.NewPoolManager(ctx, postgres.Settings{
			ConnString:          s.databaseURL,
			RetryCount:          DBConnRetryCount,
			RetryDelay:          DBConnRetryDelay,
			MigrationsDirectory: s.migrationsDir,
		})

		if err != nil {
			slog.Error("Error creating pool manager", slog.Any("err", err))
			return err
		}

		userRepo = usersrepo.NewPostgresRepository(poolManager)
		secretsRepo = secretsrepo.NewPostgresRepository(poolManager)
	}

	// Setup controllers
	userController := usercontroller.NewUserController(userRepo, jwtResolver)
	secretsController := secretscontroller.NewSecretsController(secretsRepo)

	// Setup public routes
	publicRoutes := []string{
		"/proto.EncryNestUserService/RegisterUser",
		"/proto.EncryNestUserService/AuthenticateUser",
	}

	// Setup interceptors
	auth := interceptors.NewAuthInterceptor(jwtResolver, publicRoutes)
	logger := &interceptors.LoggerInterceptor{}

	unaryInterceptors := make([]grpc.UnaryServerInterceptor, 0)
	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)

	unaryInterceptors = append(unaryInterceptors, auth.Unary())
	unaryInterceptors = append(unaryInterceptors, logger.Unary())

	streamInterceptors = append(streamInterceptors, auth.Stream())
	streamInterceptors = append(streamInterceptors, logger.Stream())

	// Setup server options
	serverOptions := make([]grpc.ServerOption, 0)
	serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	serverOptions = append(serverOptions, grpc.ChainStreamInterceptor(streamInterceptors...))

	// Setup TLS config
	if s.certsDirectory != "" || (s.certFileName != "" && s.keyFileName != "") {
		tlsLoadingSettings := keychain.TLSCertificateSettings{
			Directory: s.certsDirectory,
			CertFile:  s.certFileName,
			KeyFile:   s.keyFileName,
		}

		cert, err := keychain.LoadTLSCertificate(tlsLoadingSettings)
		if err != nil {
			return err
		}

		cred := credentials.NewServerTLSFromCert(cert)
		serverOptions = append(serverOptions, grpc.Creds(cred))
		slog.Info("TLS is enabled")
	} else {
		slog.Info("TLS is disabled")
	}

	// Setup gRPC server
	s.srv = grpc.NewServer(serverOptions...)

	// Register services
	pb.RegisterEncryNestUserServiceServer(
		s.srv,
		userservice.NewUserService(userController),
	)
	pb.RegisterEncryNestSecretsServiceServer(
		s.srv,
		secretsservice.NewSecretsService(secretsController),
	)
	// TODO: Add file service
	// TODO: Add keys service

	return nil
}

func (s *Server) Run(ctx context.Context) error {
	if s.srv == nil {
		slog.Error("Server is not initialized")
		return errors.New("server is not initialized")
	}
	// Setup listener
	listenConfig := net.ListenConfig{
		KeepAlive: -1,
	}
	listener, err := listenConfig.Listen(ctx, "tcp", s.addr+":"+s.port)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}

	// Start gRPC server
	eg.Go(func() error {
		return s.srv.Serve(listener)
	})

	slog.Info("Server started", slog.String("addr", s.addr+":"+s.port))

	// Wait for context to be done
	eg.Go(func() error {
		<-ctx.Done()
		slog.Info("Shutting down server")
		s.srv.GracefulStop()
		return nil
	})

	return eg.Wait()
}
