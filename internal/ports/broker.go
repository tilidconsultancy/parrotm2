package ports

import (
	"context"

	"github.com/google/uuid"
)

type (
	ConsumerFunc[T interface{}] func(context.Context, ConsumerContext, T)
	ConsumerContext             struct {
		RemainingRetries uint16
		Faulted          bool
	}
	Consumer[T interface{}] interface {
	}
	Producer[T interface{}] interface {
		Publish(correlationId uuid.UUID, msgs ...T) error
	}
)
