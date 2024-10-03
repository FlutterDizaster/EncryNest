package secretservice

import (
	"context"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"
	"github.com/FlutterDizaster/EncryNest/internal/models/secrets"
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
	MakeUpdate(ctx context.Context, update *secrets.Update) (string, uuid.UUID, error)

	// SubscribeUpdates subscribes for secret updates.
	// Return channel of secret updates.
	SubscribeUpdates(
		ctx context.Context,
		owner uuid.UUID,
		client uuid.UUID,
		knownVersion string,
		knownIDs []string,
	) <-chan secrets.Update
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
	secret, err := secrets.NewSecretFromProto(req.GetSecret())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid secret ID")
	}

	// Obtaining user id from ctx
	rawUserID := ctx.Value("userID")
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user ID")
	}

	// Obtaining client id from ctx
	rawClientID := ctx.Value("clientID")
	clientID, err := uuid.Parse(rawClientID.(string))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid client ID")
	}

	upd := &secrets.Update{
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
	// Obtaining owner id from ctx
	rawOwnerID := stream.Context().Value("userID")
	ownerID, err := uuid.Parse(rawOwnerID.(string))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid user ID")
	}

	// Obtaining client id from ctx
	rawClientID := stream.Context().Value("clientID")
	clientID, err := uuid.Parse(rawClientID.(string))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid client ID")
	}

	updatesChan := s.secretsController.SubscribeUpdates(
		stream.Context(),
		ownerID,
		clientID,
		req.GetKnownVersion(),
		req.GetKnownIds(),
	)

	for update := range updatesChan {
		var resp pb.Update

		resp.Secret = update.Secret.ToProto()

		switch update.Action {
		case secrets.UpdateActionCreate:
			resp.Action = pb.Action_CREATE

		case secrets.UpdateActionUpdate:
			resp.Action = pb.Action_UPDATE

		case secrets.UpdateActionDelete:
			resp.Action = pb.Action_DELETE

		default:
			// TODO: Add default case
		}

		err = stream.Send(&resp)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send update: %v", err)
		}
	}

	return nil
}
