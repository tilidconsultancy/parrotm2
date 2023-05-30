package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"pm2/internal/adapters/gRPC"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"time"

	"github.com/google/uuid"
)

type (
	ConversationServer struct {
		gRPC.UnimplementedConversationServiceServer

		conversationRepository ports.Repository[domain.Conversation]
		tenantUserRepository   ports.Repository[domain.TenantUser]
		sessionManager         ports.SessionManagerUseCase
	}
)

func NewConversationServer(conversationRepository ports.Repository[domain.Conversation],
	tenantUserRepository ports.Repository[domain.TenantUser],
	sessionManager ports.SessionManagerUseCase) gRPC.ConversationServiceServer {
	return &ConversationServer{
		conversationRepository: conversationRepository,
		sessionManager:         sessionManager,
		tenantUserRepository:   tenantUserRepository,
	}
}

func (c *ConversationServer) TakeOverConversation(rq *gRPC.TakeConversation, rw gRPC.ConversationService_TakeOverConversationServer) error {
	ctx := rw.Context()
	tuid, err := uuid.Parse(rq.TenantUserId)
	if err != nil {
		return err
	}
	tuser := c.tenantUserRepository.GetFirst(ctx, ports.GetById(tuid))
	if tuser == nil {
		return errors.New(domain.TENANT_USER_NOT_FOUND)
	}

	cvid, err := uuid.Parse(rq.ConversationId)
	if err != nil {
		return err
	}
	cv := c.conversationRepository.GetFirst(ctx, ports.GetById(cvid))
	if cv == nil || cv.TenantUser != nil {
		return errors.New(domain.CONVERSATION_NOT_FOUND)
	}
	cv.TenantUser = tuser
	c.conversationRepository.Replace(ctx, ports.GetById(cvid), cv)
	for {
		select {
		case <-ctx.Done():
			ctx := context.Background()
			cv := c.conversationRepository.GetFirst(ctx, ports.GetById(cvid))
			cv.TenantUser = nil
			c.conversationRepository.Replace(ctx, ports.GetById(cvid), cv)
			return nil
		default:
			rw.Send(rq)
			time.Sleep(time.Second * 2)
		}
	}
}

func (c *ConversationServer) GetAllConversations(
	rq *gRPC.ConversationRequest,
	rw gRPC.ConversationService_GetAllConversationsServer) error {
	ctx := rw.Context()
	j, _ := json.Marshal(rq)
	log.Println(string(j))
	id, err := uuid.Parse(rq.TenantId)
	if err != nil {
		return err
	}
	filter := map[string]interface{}{
		"tenant.id": id,
	}
	for _, f := range rq.Filters {
		filter[f.Key] = f.Value
	}
	cvs := c.conversationRepository.GetAllSkipTake(ctx, filter, rq.Skip, rq.Take)
	cr := buildConversationResponse(cvs)
	rw.Send(cr)
	j, _ = json.Marshal(cr)
	log.Println(string(j))
	s := &ports.Session{
		Id:  uuid.New(),
		Key: rq.TenantId,
	}
	c.sessionManager.CreateSession(s)
	c.sessionManager.AppendSessionEvent(func(ss *ports.Session) bool {
		return s.Id == ss.Id
	}, func(_ context.Context, i interface{}) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = e.(error)
			}
		}()
		cv := i.(*domain.Conversation)
		cr := &gRPC.ConversationResponse{
			Conversations: buildConversations(*cv),
			Count:         1,
		}
		rw.Send(cr)
		j, _ = json.Marshal(cr)
		log.Println(string(j))
		return nil
	})
	defer c.sessionManager.RemoveSessions(func(ss *ports.Session) bool {
		return ss.Id == s.Id
	})
	<-ctx.Done()
	return nil
}

func buildConversationResponse(cvs *ports.Pagination[domain.Conversation]) *gRPC.ConversationResponse {
	return &gRPC.ConversationResponse{
		Count:         cvs.Count,
		Conversations: buildConversations(cvs.Data...),
	}
}

func buildConversations(cvs ...domain.Conversation) []*gRPC.Conversation {
	r := []*gRPC.Conversation{}
	for _, c := range cvs {
		r = append(r, &gRPC.Conversation{
			Id:        c.Id.String(),
			Tenant:    buildTenant(&c.Tenant),
			User:      buildUser(&c.User),
			Status:    string(c.Status),
			CreatedAt: c.CreatedAt.String(),
			UpdatedAt: c.UpdatedAt.String(),
		})
	}
	return r
}

func buildTenant(t *domain.Tenant) *gRPC.Tenant {
	return &gRPC.Tenant{
		Id:        t.Id.String(),
		Name:      t.Name,
		Contacts:  buildContacts(t.Contacts),
		Addresses: buildAddresses(t.Addresses),
	}
}

func buildContacts(css []domain.Contact) []*gRPC.Contact {
	r := []*gRPC.Contact{}
	for _, c := range css {
		r = append(r, &gRPC.Contact{
			Label:   c.Contact,
			Contact: c.Contact,
		})
	}
	return r
}

func buildAddresses(css []domain.Address) []*gRPC.Address {
	r := []*gRPC.Address{}
	for _, c := range css {
		r = append(r, &gRPC.Address{
			Label:    c.Label,
			Zipcode:  c.Zipcode,
			Street:   c.Street,
			Number:   c.Number,
			District: c.District,
			City:     c.City,
			State:    c.State,
		})
	}
	return r
}

func buildUser(u *domain.User) *gRPC.User {
	return &gRPC.User{
		Name:         u.Name,
		Phone:        u.Phone,
		Informations: buildKeyValue(u.Informations),
	}
}

func buildKeyValue(m map[string]string) []*gRPC.KeyValue {
	r := []*gRPC.KeyValue{}
	for k, v := range m {
		r = append(r, &gRPC.KeyValue{
			Key:   k,
			Value: v,
		})
	}
	return r
}
