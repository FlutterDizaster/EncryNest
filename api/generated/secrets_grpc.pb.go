// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: api/proto/secrets.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EncryNestSecretsService_SubscribeUpdates_FullMethodName = "/proto.EncryNestSecretsService/SubscribeUpdates"
	EncryNestSecretsService_MakeUpdate_FullMethodName       = "/proto.EncryNestSecretsService/MakeUpdate"
)

// EncryNestSecretsServiceClient is the client API for EncryNestSecretsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EncryNestSecretsServiceClient interface {
	SubscribeUpdates(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Update], error)
	MakeUpdate(ctx context.Context, in *Update, opts ...grpc.CallOption) (*MakeUpdateResponse, error)
}

type encryNestSecretsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEncryNestSecretsServiceClient(cc grpc.ClientConnInterface) EncryNestSecretsServiceClient {
	return &encryNestSecretsServiceClient{cc}
}

func (c *encryNestSecretsServiceClient) SubscribeUpdates(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Update], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &EncryNestSecretsService_ServiceDesc.Streams[0], EncryNestSecretsService_SubscribeUpdates_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeRequest, Update]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EncryNestSecretsService_SubscribeUpdatesClient = grpc.ServerStreamingClient[Update]

func (c *encryNestSecretsServiceClient) MakeUpdate(ctx context.Context, in *Update, opts ...grpc.CallOption) (*MakeUpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MakeUpdateResponse)
	err := c.cc.Invoke(ctx, EncryNestSecretsService_MakeUpdate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EncryNestSecretsServiceServer is the server API for EncryNestSecretsService service.
// All implementations must embed UnimplementedEncryNestSecretsServiceServer
// for forward compatibility.
type EncryNestSecretsServiceServer interface {
	SubscribeUpdates(*SubscribeRequest, grpc.ServerStreamingServer[Update]) error
	MakeUpdate(context.Context, *Update) (*MakeUpdateResponse, error)
	mustEmbedUnimplementedEncryNestSecretsServiceServer()
}

// UnimplementedEncryNestSecretsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEncryNestSecretsServiceServer struct{}

func (UnimplementedEncryNestSecretsServiceServer) SubscribeUpdates(*SubscribeRequest, grpc.ServerStreamingServer[Update]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeUpdates not implemented")
}
func (UnimplementedEncryNestSecretsServiceServer) MakeUpdate(context.Context, *Update) (*MakeUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeUpdate not implemented")
}
func (UnimplementedEncryNestSecretsServiceServer) mustEmbedUnimplementedEncryNestSecretsServiceServer() {
}
func (UnimplementedEncryNestSecretsServiceServer) testEmbeddedByValue() {}

// UnsafeEncryNestSecretsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EncryNestSecretsServiceServer will
// result in compilation errors.
type UnsafeEncryNestSecretsServiceServer interface {
	mustEmbedUnimplementedEncryNestSecretsServiceServer()
}

func RegisterEncryNestSecretsServiceServer(s grpc.ServiceRegistrar, srv EncryNestSecretsServiceServer) {
	// If the following call pancis, it indicates UnimplementedEncryNestSecretsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EncryNestSecretsService_ServiceDesc, srv)
}

func _EncryNestSecretsService_SubscribeUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EncryNestSecretsServiceServer).SubscribeUpdates(m, &grpc.GenericServerStream[SubscribeRequest, Update]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EncryNestSecretsService_SubscribeUpdatesServer = grpc.ServerStreamingServer[Update]

func _EncryNestSecretsService_MakeUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Update)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EncryNestSecretsServiceServer).MakeUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EncryNestSecretsService_MakeUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EncryNestSecretsServiceServer).MakeUpdate(ctx, req.(*Update))
	}
	return interceptor(ctx, in, info, handler)
}

// EncryNestSecretsService_ServiceDesc is the grpc.ServiceDesc for EncryNestSecretsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EncryNestSecretsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.EncryNestSecretsService",
	HandlerType: (*EncryNestSecretsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MakeUpdate",
			Handler:    _EncryNestSecretsService_MakeUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeUpdates",
			Handler:       _EncryNestSecretsService_SubscribeUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/proto/secrets.proto",
}
