package ports

import (
	"context"
	"pm2/internal/domain"
)

type (
	NlpClient interface {
		UnrollConversation(ctx context.Context, msgs []domain.Msg) (*domain.Msg, error)
	}
)
