// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.1
// source: protobuf/urlShorter.proto

package protobuf

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

// UrlShortnerServiceClient is the client API for UrlShortnerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlShortnerServiceClient interface {
	CreateShortUrl(ctx context.Context, in *CreateShortUrlRequest, opts ...grpc.CallOption) (*CreateShortUrlResponse, error)
}

type urlShortnerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlShortnerServiceClient(cc grpc.ClientConnInterface) UrlShortnerServiceClient {
	return &urlShortnerServiceClient{cc}
}

func (c *urlShortnerServiceClient) CreateShortUrl(ctx context.Context, in *CreateShortUrlRequest, opts ...grpc.CallOption) (*CreateShortUrlResponse, error) {
	out := new(CreateShortUrlResponse)
	err := c.cc.Invoke(ctx, "/UrlShortnerService/CreateShortUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlShortnerServiceServer is the server API for UrlShortnerService service.
// All implementations must embed UnimplementedUrlShortnerServiceServer
// for forward compatibility
type UrlShortnerServiceServer interface {
	CreateShortUrl(context.Context, *CreateShortUrlRequest) (*CreateShortUrlResponse, error)
	mustEmbedUnimplementedUrlShortnerServiceServer()
}

// UnimplementedUrlShortnerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUrlShortnerServiceServer struct {
}

func (UnimplementedUrlShortnerServiceServer) CreateShortUrl(context.Context, *CreateShortUrlRequest) (*CreateShortUrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortUrl not implemented")
}
func (UnimplementedUrlShortnerServiceServer) mustEmbedUnimplementedUrlShortnerServiceServer() {}

// UnsafeUrlShortnerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlShortnerServiceServer will
// result in compilation errors.
type UnsafeUrlShortnerServiceServer interface {
	mustEmbedUnimplementedUrlShortnerServiceServer()
}

func RegisterUrlShortnerServiceServer(s grpc.ServiceRegistrar, srv UrlShortnerServiceServer) {
	s.RegisterService(&UrlShortnerService_ServiceDesc, srv)
}

func _UrlShortnerService_CreateShortUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShortUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortnerServiceServer).CreateShortUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UrlShortnerService/CreateShortUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortnerServiceServer).CreateShortUrl(ctx, req.(*CreateShortUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UrlShortnerService_ServiceDesc is the grpc.ServiceDesc for UrlShortnerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UrlShortnerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "UrlShortnerService",
	HandlerType: (*UrlShortnerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortUrl",
			Handler:    _UrlShortnerService_CreateShortUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/urlShorter.proto",
}
