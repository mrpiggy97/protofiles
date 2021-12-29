// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package formatting

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FormattingServiceClient is the client API for FormattingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FormattingServiceClient interface {
	ToCamelCase(ctx context.Context, in *FormattingRequest, opts ...grpc.CallOption) (*FormattingResponse, error)
	ToLowerCase(ctx context.Context, in *FormattingRequest, opts ...grpc.CallOption) (*FormattingResponse, error)
	ToUpperCase(ctx context.Context, in *FormattingRequest, opts ...grpc.CallOption) (*FormattingResponse, error)
}

type formattingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFormattingServiceClient(cc grpc.ClientConnInterface) FormattingServiceClient {
	return &formattingServiceClient{cc}
}

func (c *formattingServiceClient) ToCamelCase(ctx context.Context, in *FormattingRequest, opts ...grpc.CallOption) (*FormattingResponse, error) {
	out := new(FormattingResponse)
	err := c.cc.Invoke(ctx, "/strings.FormattingService/ToCamelCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *formattingServiceClient) ToLowerCase(ctx context.Context, in *FormattingRequest, opts ...grpc.CallOption) (*FormattingResponse, error) {
	out := new(FormattingResponse)
	err := c.cc.Invoke(ctx, "/strings.FormattingService/ToLowerCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *formattingServiceClient) ToUpperCase(ctx context.Context, in *FormattingRequest, opts ...grpc.CallOption) (*FormattingResponse, error) {
	out := new(FormattingResponse)
	err := c.cc.Invoke(ctx, "/strings.FormattingService/ToUpperCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FormattingServiceServer is the server API for FormattingService service.
// All implementations must embed UnimplementedFormattingServiceServer
// for forward compatibility
type FormattingServiceServer interface {
	ToCamelCase(context.Context, *FormattingRequest) (*FormattingResponse, error)
	ToLowerCase(context.Context, *FormattingRequest) (*FormattingResponse, error)
	ToUpperCase(context.Context, *FormattingRequest) (*FormattingResponse, error)
	mustEmbedUnimplementedFormattingServiceServer()
}

// UnimplementedFormattingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFormattingServiceServer struct {
}

func (UnimplementedFormattingServiceServer) ToCamelCase(context.Context, *FormattingRequest) (*FormattingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToCamelCase not implemented")
}
func (UnimplementedFormattingServiceServer) ToLowerCase(context.Context, *FormattingRequest) (*FormattingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToLowerCase not implemented")
}
func (UnimplementedFormattingServiceServer) ToUpperCase(context.Context, *FormattingRequest) (*FormattingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToUpperCase not implemented")
}
func (UnimplementedFormattingServiceServer) mustEmbedUnimplementedFormattingServiceServer() {}

// UnsafeFormattingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FormattingServiceServer will
// result in compilation errors.
type UnsafeFormattingServiceServer interface {
	mustEmbedUnimplementedFormattingServiceServer()
}

func RegisterFormattingServiceServer(s grpc.ServiceRegistrar, srv FormattingServiceServer) {
	s.RegisterService(&FormattingService_ServiceDesc, srv)
}

func _FormattingService_ToCamelCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FormattingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FormattingServiceServer).ToCamelCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strings.FormattingService/ToCamelCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FormattingServiceServer).ToCamelCase(ctx, req.(*FormattingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FormattingService_ToLowerCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FormattingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FormattingServiceServer).ToLowerCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strings.FormattingService/ToLowerCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FormattingServiceServer).ToLowerCase(ctx, req.(*FormattingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FormattingService_ToUpperCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FormattingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FormattingServiceServer).ToUpperCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strings.FormattingService/ToUpperCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FormattingServiceServer).ToUpperCase(ctx, req.(*FormattingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FormattingService_ServiceDesc is the grpc.ServiceDesc for FormattingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FormattingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "strings.FormattingService",
	HandlerType: (*FormattingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ToCamelCase",
			Handler:    _FormattingService_ToCamelCase_Handler,
		},
		{
			MethodName: "ToLowerCase",
			Handler:    _FormattingService_ToLowerCase_Handler,
		},
		{
			MethodName: "ToUpperCase",
			Handler:    _FormattingService_ToUpperCase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protofiles/Formatting.proto",
}
