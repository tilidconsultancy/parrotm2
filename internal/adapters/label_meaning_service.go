package adapters

import (
	"context"
	"pm2/internal/domain"
	"pm2/internal/ports"
)

type (
	LabelMeaningService struct {
		labelMeaningRepository ports.Repository[domain.LabelMeaning]
		nlpcli                 ports.NlpClient
	}
)

func NewLabelMeaningService(labelMeaningRepository ports.Repository[domain.LabelMeaning],
	nlpcli ports.NlpClient) ports.LabelMeaningUseCase {
	return &LabelMeaningService{
		labelMeaningRepository: labelMeaningRepository,
		nlpcli:                 nlpcli,
	}
}

func (lms *LabelMeaningService) GetPercentageMeaningsByMessages(ctx context.Context, t *domain.Tenant, msgs []domain.Msg) ([]domain.PercentageLabelMeaning, error) {
	labels := lms.labelMeaningRepository.GetAll(ctx, ports.GetByOwnerId(t.Id))
	return lms.nlpcli.EvaluateLabelMeaning(ctx, t, msgs, labels)
}
