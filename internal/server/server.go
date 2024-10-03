package server

import (
	"context"
	"log/slog"
	"net"
	"time"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"
	secretscontroller "github.com/FlutterDizaster/EncryNest/internal/server/controllers/secrets-controller"
	usercontroller "github.com/FlutterDizaster/EncryNest/internal/server/controllers/user-controller"
	jwtresolver "github.com/FlutterDizaster/EncryNest/internal/server/jwt-resolver"
	secretsrepo "github.com/FlutterDizaster/EncryNest/internal/server/repository/secrets"
	usersrepo "github.com/FlutterDizaster/EncryNest/internal/server/repository/users"
	"github.com/FlutterDizaster/EncryNest/internal/server/services/interceptors"
	secretsservice "github.com/FlutterDizaster/EncryNest/internal/server/services/secrets"
	userservice "github.com/FlutterDizaster/EncryNest/internal/server/services/users"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
)

const (
	TokenTTL = time.Hour * 24 * 31
)

type Settings struct {
	Addr      string
	Port      string
	JWTSecret string
}

type Server struct {
	addr      string
	port      string
	jwtSecret string
}

func NewServer(settings Settings) *Server {
	return &Server{
		addr:      settings.Addr,
		port:      settings.Port,
		jwtSecret: settings.JWTSecret,
	}
}

func (s *Server) Run(ctx context.Context) error {
	slog.Info("Starting server")

	// Setup JWT resolver
	jwtResolverSettings := jwtresolver.Settings{
		Secret:   s.jwtSecret,
		TokenTTL: TokenTTL,
	}
	jwtResolver := jwtresolver.New(jwtResolverSettings)

	// Setup repositories
	userRepo := usersrepo.NewInMemoryRepository()
	secretsRepo := secretsrepo.NewInMemoryRepository()

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

	// Setup gRPC server
	encryNestServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// Register services
	pb.RegisterEncryNestUserServiceServer(
		encryNestServer,
		userservice.NewUserService(userController),
	)
	pb.RegisterEncryNestSecretsServiceServer(
		encryNestServer,
		secretsservice.NewSecretsService(secretsController),
	)
	// TODO: Add file service
	// TODO: Add keys service

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
		return encryNestServer.Serve(listener)
	})

	// Wait for context to be done
	eg.Go(func() error {
		<-ctx.Done()
		encryNestServer.GracefulStop()
		return nil
	})

	return eg.Wait()
}
