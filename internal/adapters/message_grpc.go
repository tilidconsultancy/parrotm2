package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"pm2/internal/adapters/gRPC"
	"pm2/internal/domain"
	"pm2/internal/domain/events"
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

func (ms *MessageServer) ReceiveChunkedAudio(rq *gRPC.AudioChunkRequest, rw gRPC.MessageService_ReceiveChunkedAudioServer) error {
	ctx := rw.Context()
	file, err := os.Open("audio.mp3")
	if err != nil {
		return status.Error(codes.Aborted, err.Error())
	}
	defer file.Close()
	buff := make([]byte, rq.BufferSize)
	os.Remove("audio2.mp3")
	f, err := os.Create("audio2.mp3")
	if err != nil {
		return err
	}
	defer f.Close()
	for n, err := file.Read(buff); err != io.EOF; n, err = file.Read(buff) {
		select {
		case <-ctx.Done():
			return status.Error(codes.Aborted, "request cancelled")
		default:
			if err != nil {
				return status.Error(codes.Aborted, err.Error())
			}
			rw.Send(&gRPC.AudioChunkResponse{
				Buffer: buff,
				GCount: uint32(n),
			})
			f.Write(buff)
		}
	}
	return nil
}

func (ms *MessageServer) AssignConversationsMessages(rq *gRPC.AssinConversationsRequest, rw gRPC.MessageService_AssignConversationsMessagesServer) error {
	ctx := rw.Context()
	cids := []uuid.UUID{}
	for _, scvid := range rq.ConversationsId {
		cid, err := uuid.Parse(scvid)
		if err != nil {
			return status.Error(codes.Aborted, err.Error())
		}
		cids = append(cids, cid)
	}
	cvs := ms.conversationRepository.GetAll(ctx, ports.GetByIds(cids))
	msgs := []*gRPC.Message{}
	for _, cv := range cvs {
		m := buildMessages(&cv, cv.Messages[len(cv.Messages)-1])[0]
		msgs = append(msgs, m)
	}
	rw.Send(&gRPC.MessagesResponse{
		Messages: msgs,
	})
	s := &ports.Session{
		Id:   uuid.New(),
		Keys: rq.ConversationsId,
	}
	ms.sessionManager.CreateSession(s)
	ms.sessionManager.AppendSessionEvent(func(ss *ports.Session) bool {
		return ss.Id == s.Id
	}, func(ctx context.Context, i interface{}) (err error) {
		defer recoverMessage(&err)
		evt := i.(*events.MessageEvent)
		cv := ms.conversationRepository.GetFirst(ctx, ports.GetById(evt.ConversationId))
		if cv == nil {
			return errors.New(domain.CONVERSATION_NOT_FOUND)
		}
		rmsg := buildMessageResponse(cv, *evt.Message)
		rw.Send(rmsg)
		return nil
	})
	<-ctx.Done()
	return nil
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
	return buildMessages(cv, *msg)[0], nil
}

func (ms *MessageServer) GetMessagesByConversationId(
	rq *gRPC.MessagesRequest,
	rw gRPC.MessageService_GetMessagesByConversationIdServer) (err error) {
	defer recoverMessage(&err)
	ctx := rw.Context()
	j, _ := json.Marshal(rq)
	log.Println(string(j))
	cv := ms.conversationRepository.GetFirst(ctx, ports.GetById(uuid.MustParse(rq.ConversationId)))
	if cv == nil {
		return status.Error(codes.NotFound, domain.CONVERSATION_NOT_FOUND)
	}
	mr := buildMessageResponse(cv, cv.Messages...)
	s := &ports.Session{
		Id:   uuid.New(),
		Keys: []string{rq.ConversationId},
	}
	ms.sessionManager.CreateSession(s)
	ms.sessionManager.AppendSessionEvent(func(ss *ports.Session) bool {
		return s.Id == ss.Id
	}, func(_ context.Context, i interface{}) (err error) {
		defer recoverMessage(&err)
		msg := i.(*events.MessageEvent)
		msgr := buildMessageResponse(cv, *msg.Message)
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

func buildMessageResponse(cv *domain.Conversation, msgs ...domain.Msg) *gRPC.MessagesResponse {
	return &gRPC.MessagesResponse{
		Messages: buildMessages(cv, msgs...),
	}
}

func buildMessages(cv *domain.Conversation, msgs ...domain.Msg) []*gRPC.Message {
	r := []*gRPC.Message{}
	for _, m := range msgs {
		r = append(r, &gRPC.Message{
			Id:           m.Id,
			Role:         string(m.Role),
			Content:      m.Content,
			Status:       string(m.Status),
			Conversation: buildConversations(*cv)[0],
			TenantUser:   buildTenantUser(m.TenantUser),
			CreatedAt:    m.CreatedAt.String(),
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
