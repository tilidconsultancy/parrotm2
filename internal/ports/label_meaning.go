package ports

import (
	"context"
	"pm2/internal/domain"
)

type (
	LabelMeaningUseCase interface {
		GetPercentageMeaningsByMessages(
			ctx context.Context,
			t *domain.Tenant,
			msgs []domain.Msg) ([]domain.PercentageLabelMeaning, error)
	}
)
