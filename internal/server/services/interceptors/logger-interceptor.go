package interceptors

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

// LoggerInterceptor is a grpc server interceptor that logs incoming requests.
type LoggerInterceptor struct {
}

// Unary returns unary server interceptor.
//
//nolint:nonamedreturns // grpc type declaration
func (i *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		startTime := time.Now()

		resp, err = handler(ctx, req)

		slog.Info(
			"incoming unary request",
			slog.String("Method", info.FullMethod),
			slog.Int64("TimeElapsed", time.Since(startTime).Milliseconds()),
			slog.Group(
				"Response",
				slog.Bool("Success", resp != nil),
				slog.Any("Error", err),
			),
		)

		return resp, err
	}
}

// Stream returns stream server interceptor.
func (i *LoggerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		slog.Info(
			"new stream request",
			slog.String("Method", info.FullMethod),
		)

		err := handler(srv, stream)

		slog.Info(
			"finished stream request",
			slog.String("Method", info.FullMethod),
			slog.Int64("TimeElapsed", time.Since(startTime).Milliseconds()),
			slog.Any("Error", err),
		)

		return err
	}
}
