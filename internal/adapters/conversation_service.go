package adapters

import (
	"context"
	"errors"
	"pm2/internal/domain"
	"pm2/internal/domain/events"
	"pm2/internal/ports"
	"pm2/internal/ports/boundaries"
	"time"

	"github.com/google/uuid"
	"github.com/ledongthuc/goterators"
)

type (
	ConversationService struct {
		gptClient              ports.NlpClient
		conversationRepository ports.Repository[domain.Conversation]
		tenantRepository       ports.Repository[domain.Tenant]
		metaClient             ports.MetaClient
		conversationProducer   ports.Producer[events.ConversationEvent]
		messageProducer        ports.Producer[events.MessageEvent]
	}
)

func NewConversationService(gptClient ports.NlpClient,
	conversationRepository ports.Repository[domain.Conversation],
	tenantRepository ports.Repository[domain.Tenant],
	metaClient ports.MetaClient,
	conversationProducer ports.Producer[events.ConversationEvent],
	messageProducer ports.Producer[events.MessageEvent]) ports.ConversationUseCase {
	return &ConversationService{
		gptClient:              gptClient,
		conversationRepository: conversationRepository,
		tenantRepository:       tenantRepository,
		metaClient:             metaClient,
		conversationProducer:   conversationProducer,
		messageProducer:        messageProducer,
	}
}

func (cs *ConversationService) SendApplicationMessage(ctx context.Context, cv *domain.Conversation, content string) (*domain.Msg, error) {
	mid, err := cs.metaClient.SendTextMessage(ctx, &cv.Tenant, cv.User.Phone, content)
	if err != nil {
		return nil, err
	}
	msg := domain.NewMessage(
		mid,
		domain.APPLICATION,
		content,
		domain.SENT,
		cv.TenantUser)
	if err := cs.messageProducer.Publish(cv.Id, events.NewMessageEvent(cv.Id,
		cv.Id,
		msg)); err != nil {
		return nil, err
	}
	return msg, nil
}

func (cs *ConversationService) genconversation(ctx context.Context,
	cv *domain.Conversation,
	m *boundaries.Message,
	msgs []domain.Msg) (*domain.Msg, error) {
	nmsg, err := cs.gptClient.UnrollConversation(ctx, cv.Tenant.Id, msgs)
	if err != nil {
		nmsg = domain.NewMessage(
			uuid.NewString(),
			domain.APPLICATION,
			err.Error(),
			domain.ERROR,
			nil)
	}
	nmsg.Id, err = cs.metaClient.SendTextMessage(ctx, &cv.Tenant, m.From, nmsg.Content)
	if err != nil {
		return nil, err
	}
	if err := cs.messageProducer.Publish(cv.Id, events.NewMessageEvent(cv.Id,
		cv.Id,
		nmsg)); err != nil {
		return nil, err
	}
	return nmsg, nil
}

func (cs *ConversationService) UnrollConversation(ctx context.Context, msg *boundaries.IncomingMessageInput) (*domain.Msg, error) {
	c := msg.Entry[0].Changes[0]
	m := c.Value.Messages[0]
	cv := cs.conversationRepository.GetFirst(ctx, ports.GetConversationByTenantAndUser(c.Value.Metadata.PhoneNumberId, m.From))
	if cv == nil {
		t := cs.tenantRepository.GetFirst(ctx, ports.GetTenantByPhoneId(c.Value.Metadata.PhoneNumberId))
		if t == nil {
			return nil, errors.New(domain.TENANT_NOT_FOUND)
		}
		ct := c.Value.Contacts[0]
		cv = domain.NewConversation(t, ct.Profile.Name, ct.WaId)
		cid := uuid.New()
		if err := cs.conversationProducer.Publish(cid, events.NewConversationEvent(cid, cv)); err != nil {
			return nil, err
		}
	}
	_, fmi, err := goterators.Find(cv.Messages, func(item domain.Msg) bool {
		return item.Id == m.Id
	})
	if err == nil {
		return &cv.Messages[fmi+1], nil
	}
	go cs.metaClient.ReadMessage(ctx, &cv.Tenant, m.Id)
	unmsg := domain.NewMessage(
		m.Id,
		domain.USER,
		m.Text.Body,
		domain.RECEIVED,
		nil)
	if err := cs.messageProducer.Publish(cv.Id, events.NewMessageEvent(cv.Id,
		cv.Id, unmsg)); err != nil {
		return nil, err
	}
	msgs := append(cv.Messages, *unmsg)
	nmsg := unmsg
	if cv.TenantUser == nil {
		nmsg, err = cs.genconversation(ctx, cv, &m, msgs)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, *nmsg)
	}
	cv.Messages = msgs
	cv.UpdatedAt = time.Now()
	cs.conversationRepository.Replace(ctx, ports.GetById(cv.Id), cv)
	return nmsg, nil
}
