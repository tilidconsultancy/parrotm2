package ports

import (
	"context"
	"pm2/internal/domain"
)

type (
	MetaClient interface {
		SendTextMessage(ctx context.Context,
			t *domain.Tenant,
			to string,
			b string) (string, error)

		ReadMessage(ctx context.Context,
			t *domain.Tenant,
			id string) error
	}
)
