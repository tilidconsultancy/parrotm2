package ports

import (
	"context"
	"io"
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
		GetAudio(ctx context.Context,
			t *domain.Tenant,
			id string) ([]byte, error)
		SendAudioMessage(ctx context.Context,
			t *domain.Tenant,
			to string,
			id string) (string, error)
		UploadMedia(ctx context.Context,
			t *domain.Tenant,
			stream io.Reader) (string, error)
	}
)
