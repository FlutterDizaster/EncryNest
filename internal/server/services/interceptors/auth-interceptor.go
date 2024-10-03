package interceptors

import (
	"context"
	"errors"
	"slices"

	ctxvalues "github.com/FlutterDizaster/EncryNest/internal/models/ctx-values"
	jwtresolver "github.com/FlutterDizaster/EncryNest/internal/server/jwt-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Auth interceptor for gRPC server.
// Auth interceptor used to obtain user and client IDs from JWT token provided through gRPC metadata.
// And create context with user and client IDs.
// Auth interceptor must be created with NewAuthInterceptor function.
type AuthInterceptor struct {
	resolver *jwtresolver.JWTResolver
	ignored  []string
}

// NewAuthInterceptor creates new auth interceptor.
func NewAuthInterceptor(resolver *jwtresolver.JWTResolver, ignored []string) *AuthInterceptor {
	return &AuthInterceptor{
		resolver: resolver,
		ignored:  ignored,
	}
}

// Unary returns unary server interceptor.
//
//nolint:nonamedreturns // grpc type declaration
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		if slices.Contains(i.ignored, info.FullMethod) {
			return handler(ctx, req)
		}

		ctxWithIDs, err := i.createContextWithIDs(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctxWithIDs, req)
	}
}

// Stream returns stream server interceptor.
func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if slices.Contains(i.ignored, info.FullMethod) {
			return handler(srv, stream)
		}

		ctxWithIDs, err := i.createContextWithIDs(stream.Context())
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(srv, &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctxWithIDs,
		})
	}
}

func (i *AuthInterceptor) createContextWithIDs(ctx context.Context) (context.Context, error) {
	// Extracting metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md) == 0 {
		return ctx, errors.New("metadata is not provided")
	}

	// Extracting token from metadata
	token := md.Get("authorization")[0]

	// extracting claims from token
	claims, err := i.resolver.DecryptToken(token)
	if err != nil {
		return ctx, err
	}

	// Creating context with user and client IDs
	ctx = context.WithValue(ctx, ctxvalues.ContextUserID, claims.UserID)
	ctx = context.WithValue(ctx, ctxvalues.ContextClientID, claims.ClientID)

	return ctx, nil
}
