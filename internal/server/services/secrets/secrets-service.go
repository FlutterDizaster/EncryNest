package secretsservice

import (
	"context"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"
	"github.com/FlutterDizaster/EncryNest/internal/models"
	ctxvalues "github.com/FlutterDizaster/EncryNest/internal/models/ctx-values"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SecretsController is an interface that contains methods for secret management.
type SecretsController interface {
	// MakeUpdate used to create new secret in the system.
	// Return new version string and secret ID.
	// Or error.
	MakeUpdate(ctx context.Context, update *models.Update) (string, uuid.UUID, error)

	// SubscribeUpdates subscribes for secret updates.
	// Return channel of secret updates.
	SubscribeUpdates(
		ctx context.Context,
		userID uuid.UUID,
		clientID uuid.UUID,
		knownVersion string,
		knownIDs []string,
	) <-chan models.Update
}

// SecretsService represents secrets service.
// It implements pb.EncryNestSecretsServiceServer interface.
// It provides methods for secret management.
// SecretsService must be created with NewSecretsService function.
type SecretsService struct {
	pb.UnimplementedEncryNestSecretsServiceServer

	secretsController SecretsController
}

// NewSecretsService creates new instance of SecretsService.
func NewSecretsService(secretsController SecretsController) *SecretsService {
	return &SecretsService{
		secretsController: secretsController,
	}
}

// MakeUpdate used to create new secret in the system.
func (s *SecretsService) MakeUpdate(
	ctx context.Context,
	req *pb.Update,
) (*pb.MakeUpdateResponse, error) {
	resp := &pb.MakeUpdateResponse{}

	// Converting pb secret to secret model
	secret, err := models.NewSecretFromProto(req.GetSecret())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid secret ID")
	}

	// Obtaining user id from ctx
	userID, ok := ctx.Value(ctxvalues.ContextUserID).(uuid.UUID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "can't get user ID")
	}

	// Obtaining client id from ctx
	clientID, ok := ctx.Value(ctxvalues.ContextClientID).(uuid.UUID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "can't get client ID")
	}

	upd := &models.Update{
		UserID:   userID,
		ClientID: clientID,
		Secret:   secret,
	}

	version, newID, err := s.secretsController.MakeUpdate(ctx, upd)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make update: %v", err)
	}

	resp.Version = version

	if req.GetAction() == pb.Action_CREATE {
		newIDStr := newID.String()
		resp.NewId = &newIDStr
	}

	return resp, nil
}

// SubscribeUpdates subscribes for secret updates.
func (s *SecretsService) SubscribeUpdates(
	req *pb.SubscribeRequest,
	stream grpc.ServerStreamingServer[pb.Update],
) error {
	// Obtaining user id from ctx
	userID, ok := stream.Context().Value(ctxvalues.ContextUserID).(uuid.UUID)
	if !ok {
		return status.Errorf(codes.Internal, "can't get user ID")
	}

	// Obtaining client id from ctx
	clientID, ok := stream.Context().Value(ctxvalues.ContextClientID).(uuid.UUID)
	if !ok {
		return status.Errorf(codes.Internal, "can't get client ID")
	}

	updatesChan := s.secretsController.SubscribeUpdates(
		stream.Context(),
		userID,
		clientID,
		req.GetKnownVersion(),
		req.GetKnownIds(),
	)

	for update := range updatesChan {
		var resp pb.Update

		resp.Secret = update.Secret.ToProto()

		switch update.Action {
		case models.UpdateActionCreate:
			resp.Action = pb.Action_CREATE

		case models.UpdateActionUpdate:
			resp.Action = pb.Action_UPDATE

		case models.UpdateActionDelete:
			resp.Action = pb.Action_DELETE

		default:
			// TODO: Add default case
		}

		err := stream.Send(&resp)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send update: %v", err)
		}
	}

	return nil
}
