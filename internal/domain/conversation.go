package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/ledongthuc/goterators"
)

const (
	IN_PROGRESS             ConversationStatus = "inprogress"
	USER_RESPONSE_EXPIRED   ConversationStatus = "userresponseexpired"
	TENANT_RESPONSE_EXPIRED ConversationStatus = "tenantresponseexpired"
	COMPLETED               ConversationStatus = "completed"

	USER        MsgRole = "[USER]"
	SYSTEM      MsgRole = "[SYSTEM]"
	APPLICATION MsgRole = "[APPLICATION]"

	ERROR     MsgStatus = "error"
	SENT      MsgStatus = "sent"
	GENERATED MsgStatus = "generated"
	RECEIVED  MsgStatus = "received"
)

type (
	ConversationStatus string
	MsgRole            string
	Conversation       struct {
		Id         uuid.UUID
		Tenant     Tenant
		TenantUser *TenantUser
		User       User
		Messages   []Msg
		Status     ConversationStatus
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	MsgStatus string
	Msg       struct {
		Id         string
		Role       MsgRole
		Content    string
		Status     MsgStatus
		TenantUser *TenantUser
		CreatedAt  time.Time
	}
)

func NewMessage(id string, role MsgRole, c string, st MsgStatus, tu *TenantUser) *Msg {
	return &Msg{
		Id:         id,
		Role:       role,
		Content:    c,
		Status:     st,
		TenantUser: tu,
		CreatedAt:  time.Now(),
	}
}

func NewConversation(t *Tenant,
	p string,
	wid string) *Conversation {
	return &Conversation{
		Id:     uuid.New(),
		Tenant: *t,
		User: User{
			Name:  p,
			Phone: wid,
		},
		Status:    IN_PROGRESS,
		CreatedAt: time.Now(),
	}
}

func CompileMessages(msgs []Msg) string {
	return goterators.Reduce(msgs, "",
		func(
			previousValue string,
			currentValue Msg,
			_ int,
			_ []Msg) string {
			return previousValue + "#" + string(currentValue.Role) + currentValue.Content + "\n"
		})
}
