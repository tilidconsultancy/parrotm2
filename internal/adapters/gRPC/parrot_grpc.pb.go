// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gRPC

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

// ConversationServiceClient is the client API for ConversationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConversationServiceClient interface {
	GetAllConversations(ctx context.Context, in *ConversationRequest, opts ...grpc.CallOption) (ConversationService_GetAllConversationsClient, error)
	TakeOverConversation(ctx context.Context, in *ChangeConversation, opts ...grpc.CallOption) (*ChangeConversation, error)
	GiveBackConversation(ctx context.Context, in *ChangeConversation, opts ...grpc.CallOption) (*ChangeConversation, error)
}

type conversationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConversationServiceClient(cc grpc.ClientConnInterface) ConversationServiceClient {
	return &conversationServiceClient{cc}
}

func (c *conversationServiceClient) GetAllConversations(ctx context.Context, in *ConversationRequest, opts ...grpc.CallOption) (ConversationService_GetAllConversationsClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConversationService_ServiceDesc.Streams[0], "/parrot.proto.ConversationService/GetAllConversations", opts...)
	if err != nil {
		return nil, err
	}
	x := &conversationServiceGetAllConversationsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConversationService_GetAllConversationsClient interface {
	Recv() (*ConversationResponse, error)
	grpc.ClientStream
}

type conversationServiceGetAllConversationsClient struct {
	grpc.ClientStream
}

func (x *conversationServiceGetAllConversationsClient) Recv() (*ConversationResponse, error) {
	m := new(ConversationResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *conversationServiceClient) TakeOverConversation(ctx context.Context, in *ChangeConversation, opts ...grpc.CallOption) (*ChangeConversation, error) {
	out := new(ChangeConversation)
	err := c.cc.Invoke(ctx, "/parrot.proto.ConversationService/TakeOverConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *conversationServiceClient) GiveBackConversation(ctx context.Context, in *ChangeConversation, opts ...grpc.CallOption) (*ChangeConversation, error) {
	out := new(ChangeConversation)
	err := c.cc.Invoke(ctx, "/parrot.proto.ConversationService/GiveBackConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConversationServiceServer is the server API for ConversationService service.
// All implementations must embed UnimplementedConversationServiceServer
// for forward compatibility
type ConversationServiceServer interface {
	GetAllConversations(*ConversationRequest, ConversationService_GetAllConversationsServer) error
	TakeOverConversation(context.Context, *ChangeConversation) (*ChangeConversation, error)
	GiveBackConversation(context.Context, *ChangeConversation) (*ChangeConversation, error)
	mustEmbedUnimplementedConversationServiceServer()
}

// UnimplementedConversationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConversationServiceServer struct {
}

func (UnimplementedConversationServiceServer) GetAllConversations(*ConversationRequest, ConversationService_GetAllConversationsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllConversations not implemented")
}
func (UnimplementedConversationServiceServer) TakeOverConversation(context.Context, *ChangeConversation) (*ChangeConversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeOverConversation not implemented")
}
func (UnimplementedConversationServiceServer) GiveBackConversation(context.Context, *ChangeConversation) (*ChangeConversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveBackConversation not implemented")
}
func (UnimplementedConversationServiceServer) mustEmbedUnimplementedConversationServiceServer() {}

// UnsafeConversationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConversationServiceServer will
// result in compilation errors.
type UnsafeConversationServiceServer interface {
	mustEmbedUnimplementedConversationServiceServer()
}

func RegisterConversationServiceServer(s grpc.ServiceRegistrar, srv ConversationServiceServer) {
	s.RegisterService(&ConversationService_ServiceDesc, srv)
}

func _ConversationService_GetAllConversations_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConversationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConversationServiceServer).GetAllConversations(m, &conversationServiceGetAllConversationsServer{stream})
}

type ConversationService_GetAllConversationsServer interface {
	Send(*ConversationResponse) error
	grpc.ServerStream
}

type conversationServiceGetAllConversationsServer struct {
	grpc.ServerStream
}

func (x *conversationServiceGetAllConversationsServer) Send(m *ConversationResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ConversationService_TakeOverConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeConversation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConversationServiceServer).TakeOverConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/parrot.proto.ConversationService/TakeOverConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConversationServiceServer).TakeOverConversation(ctx, req.(*ChangeConversation))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConversationService_GiveBackConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeConversation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConversationServiceServer).GiveBackConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/parrot.proto.ConversationService/GiveBackConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConversationServiceServer).GiveBackConversation(ctx, req.(*ChangeConversation))
	}
	return interceptor(ctx, in, info, handler)
}

// ConversationService_ServiceDesc is the grpc.ServiceDesc for ConversationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConversationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "parrot.proto.ConversationService",
	HandlerType: (*ConversationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TakeOverConversation",
			Handler:    _ConversationService_TakeOverConversation_Handler,
		},
		{
			MethodName: "GiveBackConversation",
			Handler:    _ConversationService_GiveBackConversation_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAllConversations",
			Handler:       _ConversationService_GetAllConversations_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/parrot.proto",
}

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	GetMessagesByConversationId(ctx context.Context, in *MessagesRequest, opts ...grpc.CallOption) (MessageService_GetMessagesByConversationIdClient, error)
	AssignConversationsMessages(ctx context.Context, in *AssinConversationsRequest, opts ...grpc.CallOption) (MessageService_AssignConversationsMessagesClient, error)
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*Message, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) GetMessagesByConversationId(ctx context.Context, in *MessagesRequest, opts ...grpc.CallOption) (MessageService_GetMessagesByConversationIdClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], "/parrot.proto.MessageService/GetMessagesByConversationId", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceGetMessagesByConversationIdClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageService_GetMessagesByConversationIdClient interface {
	Recv() (*MessagesResponse, error)
	grpc.ClientStream
}

type messageServiceGetMessagesByConversationIdClient struct {
	grpc.ClientStream
}

func (x *messageServiceGetMessagesByConversationIdClient) Recv() (*MessagesResponse, error) {
	m := new(MessagesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *messageServiceClient) AssignConversationsMessages(ctx context.Context, in *AssinConversationsRequest, opts ...grpc.CallOption) (MessageService_AssignConversationsMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[1], "/parrot.proto.MessageService/AssignConversationsMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceAssignConversationsMessagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageService_AssignConversationsMessagesClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type messageServiceAssignConversationsMessagesClient struct {
	grpc.ClientStream
}

func (x *messageServiceAssignConversationsMessagesClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *messageServiceClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/parrot.proto.MessageService/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	GetMessagesByConversationId(*MessagesRequest, MessageService_GetMessagesByConversationIdServer) error
	AssignConversationsMessages(*AssinConversationsRequest, MessageService_AssignConversationsMessagesServer) error
	SendMessage(context.Context, *SendMessageRequest) (*Message, error)
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) GetMessagesByConversationId(*MessagesRequest, MessageService_GetMessagesByConversationIdServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMessagesByConversationId not implemented")
}
func (UnimplementedMessageServiceServer) AssignConversationsMessages(*AssinConversationsRequest, MessageService_AssignConversationsMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method AssignConversationsMessages not implemented")
}
func (UnimplementedMessageServiceServer) SendMessage(context.Context, *SendMessageRequest) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_GetMessagesByConversationId_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MessagesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageServiceServer).GetMessagesByConversationId(m, &messageServiceGetMessagesByConversationIdServer{stream})
}

type MessageService_GetMessagesByConversationIdServer interface {
	Send(*MessagesResponse) error
	grpc.ServerStream
}

type messageServiceGetMessagesByConversationIdServer struct {
	grpc.ServerStream
}

func (x *messageServiceGetMessagesByConversationIdServer) Send(m *MessagesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MessageService_AssignConversationsMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AssinConversationsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageServiceServer).AssignConversationsMessages(m, &messageServiceAssignConversationsMessagesServer{stream})
}

type MessageService_AssignConversationsMessagesServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type messageServiceAssignConversationsMessagesServer struct {
	grpc.ServerStream
}

func (x *messageServiceAssignConversationsMessagesServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _MessageService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/parrot.proto.MessageService/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "parrot.proto.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _MessageService_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetMessagesByConversationId",
			Handler:       _MessageService_GetMessagesByConversationId_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AssignConversationsMessages",
			Handler:       _MessageService_AssignConversationsMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/parrot.proto",
}
