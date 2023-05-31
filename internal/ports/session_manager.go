package ports

import (
	"context"

	"github.com/google/uuid"
)

type (
	Session struct {
		Id   uuid.UUID
		Key  string
		Data map[string]string
	}
	PredicateSessionFunc  func(*Session) bool
	EventSessionFunc      func(context.Context, interface{}) error
	SessionManagerUseCase interface {
		CreateSession(*Session)
		RemoveSessions(PredicateSessionFunc)
		AppendSessionEvent(PredicateSessionFunc, EventSessionFunc)
		InvokeSessionEvents(PredicateSessionFunc, context.Context, interface{}) error
	}
)
