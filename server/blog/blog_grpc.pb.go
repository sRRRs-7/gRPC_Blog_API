// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: protoc/blog.proto

package blog

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

// BlogApiClient is the client API for BlogApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlogApiClient interface {
	CreateBlog(ctx context.Context, in *CreateBlogReq, opts ...grpc.CallOption) (*CreateBlogRes, error)
	FindBlog(ctx context.Context, in *FindBlogReq, opts ...grpc.CallOption) (BlogApi_FindBlogClient, error)
}

type blogApiClient struct {
	cc grpc.ClientConnInterface
}

func NewBlogApiClient(cc grpc.ClientConnInterface) BlogApiClient {
	return &blogApiClient{cc}
}

func (c *blogApiClient) CreateBlog(ctx context.Context, in *CreateBlogReq, opts ...grpc.CallOption) (*CreateBlogRes, error) {
	out := new(CreateBlogRes)
	err := c.cc.Invoke(ctx, "/blog.BlogApi/CreateBlog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blogApiClient) FindBlog(ctx context.Context, in *FindBlogReq, opts ...grpc.CallOption) (BlogApi_FindBlogClient, error) {
	stream, err := c.cc.NewStream(ctx, &BlogApi_ServiceDesc.Streams[0], "/blog.BlogApi/FindBlog", opts...)
	if err != nil {
		return nil, err
	}
	x := &blogApiFindBlogClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BlogApi_FindBlogClient interface {
	Recv() (*FindBlogRes, error)
	grpc.ClientStream
}

type blogApiFindBlogClient struct {
	grpc.ClientStream
}

func (x *blogApiFindBlogClient) Recv() (*FindBlogRes, error) {
	m := new(FindBlogRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BlogApiServer is the server API for BlogApi service.
// All implementations should embed UnimplementedBlogApiServer
// for forward compatibility
type BlogApiServer interface {
	CreateBlog(context.Context, *CreateBlogReq) (*CreateBlogRes, error)
	FindBlog(*FindBlogReq, BlogApi_FindBlogServer) error
}

// UnimplementedBlogApiServer should be embedded to have forward compatible implementations.
type UnimplementedBlogApiServer struct {
}

func (UnimplementedBlogApiServer) CreateBlog(context.Context, *CreateBlogReq) (*CreateBlogRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBlog not implemented")
}
func (UnimplementedBlogApiServer) FindBlog(*FindBlogReq, BlogApi_FindBlogServer) error {
	return status.Errorf(codes.Unimplemented, "method FindBlog not implemented")
}

// UnsafeBlogApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlogApiServer will
// result in compilation errors.
type UnsafeBlogApiServer interface {
	mustEmbedUnimplementedBlogApiServer()
}

func RegisterBlogApiServer(s grpc.ServiceRegistrar, srv BlogApiServer) {
	s.RegisterService(&BlogApi_ServiceDesc, srv)
}

func _BlogApi_CreateBlog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBlogReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlogApiServer).CreateBlog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blog.BlogApi/CreateBlog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlogApiServer).CreateBlog(ctx, req.(*CreateBlogReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlogApi_FindBlog_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FindBlogReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BlogApiServer).FindBlog(m, &blogApiFindBlogServer{stream})
}

type BlogApi_FindBlogServer interface {
	Send(*FindBlogRes) error
	grpc.ServerStream
}

type blogApiFindBlogServer struct {
	grpc.ServerStream
}

func (x *blogApiFindBlogServer) Send(m *FindBlogRes) error {
	return x.ServerStream.SendMsg(m)
}

// BlogApi_ServiceDesc is the grpc.ServiceDesc for BlogApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlogApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "blog.BlogApi",
	HandlerType: (*BlogApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBlog",
			Handler:    _BlogApi_CreateBlog_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindBlog",
			Handler:       _BlogApi_FindBlog_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protoc/blog.proto",
}
