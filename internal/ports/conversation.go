package ports

import (
	"context"
	"pm2/internal/domain"
	"pm2/internal/ports/boundaries"
)

type (
	ConversationUseCase interface {
		UnrollConversation(ctx context.Context, msg *boundaries.IncomingMessageInput) (*domain.Msg, error)
	}
)
