package adapters

import (
	"context"
	"pm2/internal/ports"

	"github.com/ledongthuc/goterators"
)

type (
	SessionEvents struct {
		session *ports.Session
		events  []ports.EventSessionFunc
	}
	SessionManagerService struct {
		sessions []SessionEvents
	}
)

func NewSessionManagerService() ports.SessionManagerUseCase {
	return &SessionManagerService{}
}

func (sms *SessionManagerService) CreateSession(s *ports.Session) {
	sms.sessions = append(sms.sessions, SessionEvents{
		session: s,
	})
}

func (sms *SessionManagerService) RemoveSessions(p ports.PredicateSessionFunc) {
	sms.sessions = goterators.Filter(sms.sessions, func(s SessionEvents) bool { return !p(s.session) })
}

func (sms *SessionManagerService) AppendSessionEvent(p ports.PredicateSessionFunc, e ports.EventSessionFunc) {
	for i := 0; i < len(sms.sessions); i++ {
		if p(sms.sessions[i].session) {
			sms.sessions[i].events = append(sms.sessions[i].events, e)
		}
	}
}

func (sms *SessionManagerService) InvokeSessionEvents(p ports.PredicateSessionFunc, ctx context.Context, evt interface{}) error {
	ss := goterators.Filter(sms.sessions, func(s SessionEvents) bool { return p(s.session) })
	for _, s := range ss {
		for _, e := range s.events {
			if err := e(ctx, evt); err != nil {
				return err
			}
		}
	}
	return nil
}
