package ports

import (
	"context"
	"pm2/internal/domain"

	"github.com/google/uuid"
)

type (
	NlpClient interface {
		UnrollConversation(ctx context.Context, tenantId uuid.UUID, msgs []domain.Msg) (*domain.Msg, error)
	}
)
