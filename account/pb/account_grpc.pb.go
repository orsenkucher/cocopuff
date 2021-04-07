// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// AccountingClient is the client API for Accounting service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountingClient interface {
	PostAccount(ctx context.Context, in *PostAccountRequest, opts ...grpc.CallOption) (*PostAccountResponse, error)
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	GetAccounts(ctx context.Context, in *GetAccountsRequest, opts ...grpc.CallOption) (*GetAccountsResponse, error)
}

type accountingClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountingClient(cc grpc.ClientConnInterface) AccountingClient {
	return &accountingClient{cc}
}

func (c *accountingClient) PostAccount(ctx context.Context, in *PostAccountRequest, opts ...grpc.CallOption) (*PostAccountResponse, error) {
	out := new(PostAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Accounting/PostAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Accounting/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingClient) GetAccounts(ctx context.Context, in *GetAccountsRequest, opts ...grpc.CallOption) (*GetAccountsResponse, error) {
	out := new(GetAccountsResponse)
	err := c.cc.Invoke(ctx, "/pb.Accounting/GetAccounts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountingServer is the server API for Accounting service.
// All implementations must embed UnimplementedAccountingServer
// for forward compatibility
type AccountingServer interface {
	PostAccount(context.Context, *PostAccountRequest) (*PostAccountResponse, error)
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	GetAccounts(context.Context, *GetAccountsRequest) (*GetAccountsResponse, error)
	mustEmbedUnimplementedAccountingServer()
}

// UnimplementedAccountingServer must be embedded to have forward compatible implementations.
type UnimplementedAccountingServer struct {
}

func (UnimplementedAccountingServer) PostAccount(context.Context, *PostAccountRequest) (*PostAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostAccount not implemented")
}
func (UnimplementedAccountingServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedAccountingServer) GetAccounts(context.Context, *GetAccountsRequest) (*GetAccountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccounts not implemented")
}
func (UnimplementedAccountingServer) mustEmbedUnimplementedAccountingServer() {}

// UnsafeAccountingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountingServer will
// result in compilation errors.
type UnsafeAccountingServer interface {
	mustEmbedUnimplementedAccountingServer()
}

func RegisterAccountingServer(s grpc.ServiceRegistrar, srv AccountingServer) {
	s.RegisterService(&Accounting_ServiceDesc, srv)
}

func _Accounting_PostAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServer).PostAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Accounting/PostAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServer).PostAccount(ctx, req.(*PostAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Accounting_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Accounting/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Accounting_GetAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServer).GetAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Accounting/GetAccounts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServer).GetAccounts(ctx, req.(*GetAccountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Accounting_ServiceDesc is the grpc.ServiceDesc for Accounting service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Accounting_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Accounting",
	HandlerType: (*AccountingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostAccount",
			Handler:    _Accounting_PostAccount_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _Accounting_GetAccount_Handler,
		},
		{
			MethodName: "GetAccounts",
			Handler:    _Accounting_GetAccounts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
