package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"pm2/internal/adapters/gRPC"
	"pm2/internal/domain"
	"pm2/internal/ports"

	"github.com/google/uuid"
)

type (
	MessageServer struct {
		gRPC.UnimplementedMessageServiceServer
		conversationRepository ports.Repository[domain.Conversation]
		sessionManager         ports.SessionManagerUseCase
	}
)

func NewMessageServer(conversationRepository ports.Repository[domain.Conversation],
	sessionManager ports.SessionManagerUseCase) gRPC.MessageServiceServer {
	return &MessageServer{
		conversationRepository: conversationRepository,
		sessionManager:         sessionManager,
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

func (ms *MessageServer) GetMessagesByConversationId(
	rq *gRPC.MessagesRequest,
	rw gRPC.MessageService_GetMessagesByConversationIdServer) (err error) {
	defer recoverMessage(&err)
	ctx := rw.Context()
	j, _ := json.Marshal(rq)
	log.Println(string(j))
	cv := ms.conversationRepository.GetFirst(ctx, ports.GetConversationById(uuid.MustParse(rq.ConversationId)))
	if cv == nil {
		return errors.New(domain.CONVERSATION_NOT_FOUND)
	}
	mr := buildMessageResponse(cv.Messages...)
	s := &ports.Session{
		Id:  uuid.New(),
		Key: rq.ConversationId,
	}
	ms.sessionManager.CreateSession(s)
	ms.sessionManager.AppendSessionEvent(func(ss *ports.Session) bool {
		return s.Id == ss.Id
	}, func(ctx context.Context, i interface{}) (err error) {
		defer recoverMessage(&err)
		msg := i.([]domain.Msg)
		msgr := buildMessageResponse(msg...)
		rw.Send(msgr)
		j, _ = json.Marshal(msgr)
		log.Println(string(j))
		return nil
	})
	defer ms.sessionManager.RemoveSessions(func(ss *ports.Session) bool {
		return ss.Id == s.Id
	})
	rw.Send(mr)
	j, _ = json.Marshal(mr)
	log.Println(string(j))
	<-ctx.Done()
	return nil
}

func buildMessageResponse(msgs ...domain.Msg) *gRPC.MessagesResponse {
	return &gRPC.MessagesResponse{
		Messages: buildMessages(msgs),
	}
}

func buildMessages(msgs []domain.Msg) []*gRPC.Message {
	r := []*gRPC.Message{}
	for _, m := range msgs {
		r = append(r, &gRPC.Message{
			Id:         m.Id,
			Role:       string(m.Role),
			Content:    m.Content,
			Status:     m.Status,
			TenantUser: buildTenantUser(m.TenantUser),
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
