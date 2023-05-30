package ports

import (
	"context"
	"pm2/internal/domain"
	"pm2/internal/ports/boundaries"
)

type (
	ConversationUseCase interface {
		UnrollConversation(ctx context.Context, msg *boundaries.IncomingMessageInput) (*domain.Msg, error)
		SendApplicationMessage(ctx context.Context, cv *domain.Conversation, content string) (*domain.Msg, error)
	}
)
