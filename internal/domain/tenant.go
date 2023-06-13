package domain

import (
	"errors"

	"github.com/google/uuid"
)

type (
	AccountSettings struct {
		PhoneId     string
		Token       string
		NlpToken    string
		Voice       string
		RateVoice   int
		MainContext string
	}
	Tenant struct {
		Id              uuid.UUID
		Name            string
		Contacts        []Contact
		Addresses       []Address
		AccountSettings AccountSettings
	}
	TenantUser struct {
		Id        uuid.UUID
		TenantId  uuid.UUID
		Name      string
		Contacts  []Contact
		Addresses []Address
	}
)

func NewTenant(name string,
	contacts []Contact,
	addresses []Address) *Tenant {
	return &Tenant{
		Id:        uuid.New(),
		Name:      name,
		Contacts:  contacts,
		Addresses: addresses,
	}
}

func (t *TenantUser) Validate() error {
	for _, c := range t.Contacts {
		if err := c.validate(); err != nil {
			return err
		}
	}
	for _, a := range t.Addresses {
		if err := a.validate(); err != nil {
			return err
		}
	}
	if t.Name == "" || t.TenantId == uuid.Nil {
		return errors.New(TENANT_INVALID)
	}
	return nil
}

func (t *Tenant) Validate() error {
	for _, c := range t.Contacts {
		if err := c.validate(); err != nil {
			return err
		}
	}
	for _, a := range t.Addresses {
		if err := a.validate(); err != nil {
			return err
		}
	}
	if t.Name == "" {
		return errors.New(TENANT_INVALID)
	}
	return nil
}
