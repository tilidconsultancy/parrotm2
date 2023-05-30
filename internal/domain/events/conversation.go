package events

import (
	"pm2/internal/domain"

	"github.com/google/uuid"
)

type (
	ConversationEvent struct {
		CorrelationId uuid.UUID
		Conversation  *domain.Conversation
	}
	MessageEvent struct {
		CorrelationId  uuid.UUID
		ConversationId uuid.UUID
		Message        *domain.Msg
	}
)

func NewConversationEvent(cid uuid.UUID, cv *domain.Conversation) ConversationEvent {
	return ConversationEvent{
		CorrelationId: cid,
		Conversation:  cv,
	}
}

func NewMessageEvent(cid uuid.UUID,
	cvid uuid.UUID,
	msg *domain.Msg) MessageEvent {
	return MessageEvent{
		CorrelationId:  cid,
		ConversationId: cvid,
		Message:        msg,
	}
}
