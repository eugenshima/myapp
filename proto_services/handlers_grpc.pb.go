// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: handlers.proto

package myapp

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

// PersonHandlerClient is the client API for PersonHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PersonHandlerClient interface {
	GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error)
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
}

type personHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewPersonHandlerClient(cc grpc.ClientConnInterface) PersonHandlerClient {
	return &personHandlerClient{cc}
}

func (c *personHandlerClient) GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error) {
	out := new(GetByIDResponse)
	err := c.cc.Invoke(ctx, "/personHandler/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personHandlerClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, "/personHandler/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personHandlerClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/personHandler/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personHandlerClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/personHandler/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personHandlerClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/personHandler/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PersonHandlerServer is the server API for PersonHandler service.
// All implementations must embed UnimplementedPersonHandlerServer
// for forward compatibility
type PersonHandlerServer interface {
	GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error)
	GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	mustEmbedUnimplementedPersonHandlerServer()
}

// UnimplementedPersonHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedPersonHandlerServer struct {
}

func (UnimplementedPersonHandlerServer) GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedPersonHandlerServer) GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedPersonHandlerServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPersonHandlerServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPersonHandlerServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedPersonHandlerServer) mustEmbedUnimplementedPersonHandlerServer() {}

// UnsafePersonHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PersonHandlerServer will
// result in compilation errors.
type UnsafePersonHandlerServer interface {
	mustEmbedUnimplementedPersonHandlerServer()
}

func RegisterPersonHandlerServer(s grpc.ServiceRegistrar, srv PersonHandlerServer) {
	s.RegisterService(&PersonHandler_ServiceDesc, srv)
}

func _PersonHandler_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonHandlerServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personHandler/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonHandlerServer).GetByID(ctx, req.(*GetByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonHandler_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonHandlerServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personHandler/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonHandlerServer).GetAll(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonHandler_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonHandlerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personHandler/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonHandlerServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonHandler_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonHandlerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personHandler/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonHandlerServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonHandler_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonHandlerServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personHandler/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonHandlerServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PersonHandler_ServiceDesc is the grpc.ServiceDesc for PersonHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PersonHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personHandler",
	HandlerType: (*PersonHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByID",
			Handler:    _PersonHandler_GetByID_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _PersonHandler_GetAll_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _PersonHandler_Delete_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _PersonHandler_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _PersonHandler_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "handlers.proto",
}

// UserHandlerClient is the client API for UserHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserHandlerClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*SignUpResponse, error)
	GetAll(ctx context.Context, in *UserGetAllRequest, opts ...grpc.CallOption) (*UserGetAllResponse, error)
	RefreshTokenPair(ctx context.Context, in *RefreshTokenPairRequest, opts ...grpc.CallOption) (*RefreshTokenPairResponse, error)
	Delete(ctx context.Context, in *UserDeleteRequest, opts ...grpc.CallOption) (*UserDeleteResponse, error)
}

type userHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewUserHandlerClient(cc grpc.ClientConnInterface) UserHandlerClient {
	return &userHandlerClient{cc}
}

func (c *userHandlerClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/userHandler/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userHandlerClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*SignUpResponse, error) {
	out := new(SignUpResponse)
	err := c.cc.Invoke(ctx, "/userHandler/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userHandlerClient) GetAll(ctx context.Context, in *UserGetAllRequest, opts ...grpc.CallOption) (*UserGetAllResponse, error) {
	out := new(UserGetAllResponse)
	err := c.cc.Invoke(ctx, "/userHandler/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userHandlerClient) RefreshTokenPair(ctx context.Context, in *RefreshTokenPairRequest, opts ...grpc.CallOption) (*RefreshTokenPairResponse, error) {
	out := new(RefreshTokenPairResponse)
	err := c.cc.Invoke(ctx, "/userHandler/RefreshTokenPair", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userHandlerClient) Delete(ctx context.Context, in *UserDeleteRequest, opts ...grpc.CallOption) (*UserDeleteResponse, error) {
	out := new(UserDeleteResponse)
	err := c.cc.Invoke(ctx, "/userHandler/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserHandlerServer is the server API for UserHandler service.
// All implementations must embed UnimplementedUserHandlerServer
// for forward compatibility
type UserHandlerServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	SignUp(context.Context, *SignUpRequest) (*SignUpResponse, error)
	GetAll(context.Context, *UserGetAllRequest) (*UserGetAllResponse, error)
	RefreshTokenPair(context.Context, *RefreshTokenPairRequest) (*RefreshTokenPairResponse, error)
	Delete(context.Context, *UserDeleteRequest) (*UserDeleteResponse, error)
	mustEmbedUnimplementedUserHandlerServer()
}

// UnimplementedUserHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedUserHandlerServer struct {
}

func (UnimplementedUserHandlerServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserHandlerServer) SignUp(context.Context, *SignUpRequest) (*SignUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedUserHandlerServer) GetAll(context.Context, *UserGetAllRequest) (*UserGetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedUserHandlerServer) RefreshTokenPair(context.Context, *RefreshTokenPairRequest) (*RefreshTokenPairResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshTokenPair not implemented")
}
func (UnimplementedUserHandlerServer) Delete(context.Context, *UserDeleteRequest) (*UserDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedUserHandlerServer) mustEmbedUnimplementedUserHandlerServer() {}

// UnsafeUserHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserHandlerServer will
// result in compilation errors.
type UnsafeUserHandlerServer interface {
	mustEmbedUnimplementedUserHandlerServer()
}

func RegisterUserHandlerServer(s grpc.ServiceRegistrar, srv UserHandlerServer) {
	s.RegisterService(&UserHandler_ServiceDesc, srv)
}

func _UserHandler_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userHandler/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserHandler_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userHandler/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserHandler_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserGetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userHandler/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).GetAll(ctx, req.(*UserGetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserHandler_RefreshTokenPair_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenPairRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).RefreshTokenPair(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userHandler/RefreshTokenPair",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).RefreshTokenPair(ctx, req.(*RefreshTokenPairRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserHandler_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userHandler/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).Delete(ctx, req.(*UserDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserHandler_ServiceDesc is the grpc.ServiceDesc for UserHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "userHandler",
	HandlerType: (*UserHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _UserHandler_Login_Handler,
		},
		{
			MethodName: "SignUp",
			Handler:    _UserHandler_SignUp_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _UserHandler_GetAll_Handler,
		},
		{
			MethodName: "RefreshTokenPair",
			Handler:    _UserHandler_RefreshTokenPair_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserHandler_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "handlers.proto",
}
