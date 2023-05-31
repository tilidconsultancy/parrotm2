package adapters

import (
	"context"
	"errors"
	"pm2/internal/adapters/gRPC"
	"pm2/internal/domain"
	"pm2/internal/ports"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	MessageServer struct {
		gRPC.UnimplementedMessageServiceServer
		conversationRepository ports.Repository[domain.Conversation]
		sessionManager         ports.SessionManagerUseCase
		conversationUseCase    ports.ConversationUseCase
	}
)

func NewMessageServer(conversationRepository ports.Repository[domain.Conversation],
	sessionManager ports.SessionManagerUseCase,
	conversationUseCase ports.ConversationUseCase) gRPC.MessageServiceServer {
	return &MessageServer{
		conversationRepository: conversationRepository,
		sessionManager:         sessionManager,
		conversationUseCase:    conversationUseCase,
	}
}

func recoverMessage(err *error) {
	if e := recover(); e != nil {
		switch ee := e.(type) {
		case error:
			*err = ee
		case string:
			*err = errors.New(ee)
		default:
			*err = errors.New("i really don't know")
		}
	}
}

func (ms *MessageServer) SendMessage(ctx context.Context, rq *gRPC.SendMessageRequest) (*gRPC.Message, error) {
	cvid, err := uuid.Parse(rq.ConversationId)
	if err != nil {
		return nil, err
	}
	cv := ms.conversationRepository.GetFirst(ctx, ports.GetById(cvid))
	if cv == nil {
		return nil, status.Error(codes.NotFound, domain.CONVERSATION_NOT_FOUND)
	}
	if cv.TenantUser == nil {
		return nil, status.Error(codes.NotFound, domain.TENANT_USER_NOT_FOUND)
	}
	msg, err := ms.conversationUseCase.SendApplicationMessage(ctx, cv, rq.Content)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	cv.Messages = append(cv.Messages, *msg)
	ms.conversationRepository.Replace(ctx, ports.GetById(cv.Id), cv)
	return buildMessages(cv.Id, *msg)[0], nil
}

func (ms *MessageServer) GetMessagesByConversationId(
	rq *gRPC.MessagesRequest,
	rw gRPC.MessageService_GetMessagesByConversationIdServer) (err error) {
	defer recoverMessage(&err)
	ctx := rw.Context()
	ids := []uuid.UUID{}
	for _, v := range rq.ConversationIds {
		id, err := uuid.Parse(v)
		if err != nil {
			return status.Error(codes.Aborted, err.Error())
		}
		ids = append(ids, id)
	}
	cv := ms.conversationRepository.GetFirst(ctx, ports.GetByIds(ids))
	if cv == nil {
		return status.Error(codes.NotFound, domain.CONVERSATION_NOT_FOUND)
	}
	mr := buildMessageResponse(cv.Id, cv.Messages...)
	s := &ports.Session{
		Id:   uuid.New(),
		Keys: rq.ConversationIds,
	}
	ms.sessionManager.CreateSession(s)
	ms.sessionManager.AppendSessionEvent(func(ss *ports.Session) bool {
		return s.Id == ss.Id
	}, func(_ context.Context, i interface{}) (err error) {
		defer recoverMessage(&err)
		msg := i.(*domain.Msg)
		msgr := buildMessageResponse(cv.Id, *msg)
		rw.Send(msgr)
		return nil
	})
	defer ms.sessionManager.RemoveSessions(func(ss *ports.Session) bool {
		return ss.Id == s.Id
	})
	rw.Send(mr)
	<-ctx.Done()
	return nil
}

func buildMessageResponse(cvid uuid.UUID, msgs ...domain.Msg) *gRPC.MessagesResponse {
	return &gRPC.MessagesResponse{
		Messages: buildMessages(cvid, msgs...),
	}
}

func buildMessages(cvid uuid.UUID, msgs ...domain.Msg) []*gRPC.Message {
	r := []*gRPC.Message{}
	for _, m := range msgs {
		r = append(r, &gRPC.Message{
			Id:             m.Id,
			Role:           string(m.Role),
			Content:        m.Content,
			Status:         string(m.Status),
			ConversationId: cvid.String(),
			TenantUser:     buildTenantUser(m.TenantUser),
			CreatedAt:      m.CreatedAt.String(),
		})
	}
	return r
}

func buildTenantUser(t *domain.TenantUser) *gRPC.TennantUser {
	if t == nil {
		return nil
	}
	return &gRPC.TennantUser{
		Id:        t.Id.String(),
		TenantId:  t.TenantId.String(),
		Name:      t.Name,
		Contacts:  buildContacts(t.Contacts),
		Addresses: buildAddresses(t.Addresses),
	}
}
