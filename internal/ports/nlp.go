package ports

import (
	"context"
	"pm2/internal/domain"
)

type (
	NlpClient interface {
		UnrollConversation(ctx context.Context, tenant *domain.Tenant, msgs []domain.Msg) (*domain.Msg, error)
		EvaluateLabelMeaning(ctx context.Context, tenant *domain.Tenant, msgs []domain.Msg, labels []domain.LabelMeaning) ([]domain.PercentageLabelMeaning, error)
	}
)
