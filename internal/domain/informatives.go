package domain

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"github.com/paemuri/brdoc"
)

const (
	EMAIL  ContactLabel = "email"
	MOBILE ContactLabel = "mobile"
	PHONE  ContactLabel = "phone"
)

type (
	ContactLabel string
	Contact      struct {
		Label   ContactLabel
		Contact string
	}
	Address struct {
		Label    string
		Zipcode  string
		Street   string
		Number   string
		District string
		City     string
		State    string
	}
)

func (c *Contact) Format() error {
	switch c.Label {
	case PHONE, MOBILE:
		r, _ := regexp.Compile("[^0-9]+")
		c.Contact = r.ReplaceAllString(c.Contact, "")
	}
	return c.validate()
}

func formatStrings(s ...*string) {
	for _, ss := range s {
		sss := strings.Trim(*ss, "")
		sss = strings.ToUpper(sss)
		*ss = sss
	}
}

func (a *Address) Format() error {
	r, _ := regexp.Compile("[^0-9]+")
	formatStrings(&a.City,
		&a.District,
		&a.Number,
		&a.State,
		&a.Street,
		&a.Zipcode)
	r.ReplaceAllString(a.Zipcode, "")
	return a.validate()
}

func (c *Contact) validate() error {
	if c.Contact == "" {
		return errors.New(CONTACT_INVALID)
	}
	switch c.Label {
	case EMAIL:
		_, err := mail.ParseAddress(c.Contact)
		if err != nil {
			return err
		}
	case MOBILE, PHONE:
		r := len([]rune(c.Contact))
		if r != 11 && r != 10 {
			return errors.New(CONTACT_INVALID)
		}
	}
	return nil
}

func (a *Address) validate() error {
	if !brdoc.IsCEP(a.Zipcode) ||
		a.City == "" ||
		a.District == "" ||
		a.Label == "" ||
		a.Number == "" ||
		a.State == "" ||
		a.Street == "" ||
		a.Zipcode == "" {
		return errors.New(ADDRESS_INVALID)
	}
	return nil
}
