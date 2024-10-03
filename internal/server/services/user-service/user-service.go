package userservice

import (
	"context"
	"errors"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"
	"github.com/FlutterDizaster/EncryNest/internal/models"
	sharederrors "github.com/FlutterDizaster/EncryNest/internal/shared-errors"
	"github.com/FlutterDizaster/EncryNest/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserController contains methods for user authentication and registration.
type UserController interface {
	// RegisterUser registers new user in the system and returns JWT token string.
	RegisterUser(ctx context.Context, user *models.UserCredentials) (string, error)

	// AuthUser authenticates user and returns JWT token string.
	// User can provide only username and password.
	// Email will be ignored.
	AuthUser(ctx context.Context, user *models.UserCredentials) (string, error)
}

// UserService used for user authentication and registration through gRPC.
// UserService must be created with NewUserService function.
type UserService struct {
	pb.UnimplementedEncryNestUserServiceServer

	userController UserController
}

// NewUserService creates new UserService instance with provided userController.
func NewUserService(userController UserController) *UserService {
	return &UserService{
		userController: userController,
	}
}

// RegisterUser registers new user in the system and returns JWT token string.
// RegisterUser implements EncryNestUserServiceServer.RegisterUser.
func (s *UserService) RegisterUser(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	username := req.GetUsername()
	email := req.GetEmail()
	password := req.GetPassword()

	err := validator.ValidateUsername(username)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = validator.ValidateEmail(email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = validator.ValidatePassword(password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &models.UserCredentials{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}

	token, err := s.userController.RegisterUser(ctx, user)

	if err != nil {
		if errors.Is(err, sharederrors.ErrEmailTaken) {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"user with %s email is already registered",
				req.GetEmail(),
			)
		}
		if errors.Is(err, sharederrors.ErrUsernameTaken) {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"user with %s username is already registered",
				req.GetUsername(),
			)
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	resp := &pb.RegisterUserResponse{
		Token: token,
	}

	return resp, nil
}

// AuthenticateUser used to authenticating existed user in the system and returns JWT token string.
// AuthenticateUser implements EncryNestUserServiceServer.AuthenticateUser.
func (s *UserService) AuthenticateUser(
	ctx context.Context,
	req *pb.AuthenticateUserRequest,
) (*pb.AuthenticateUserResponse, error) {
	user := &models.UserCredentials{
		Username:     req.GetUsername(),
		PasswordHash: req.GetPassword(),
	}

	token, err := s.userController.AuthUser(ctx, user)

	if err != nil {
		if errors.Is(err, sharederrors.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with given credentials not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to authenticate user: %v", err)
	}

	resp := &pb.AuthenticateUserResponse{
		Token: token,
	}

	return resp, nil
}
