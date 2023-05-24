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
		Message        []domain.Msg
	}
)
