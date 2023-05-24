package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"pm2/internal/ports/boundaries"
	"time"

	"github.com/spf13/viper"
)

type (
	MetaSettings struct {
		Retries uint8
		Host    string
		Timeout time.Duration
	}
	MetaHttpClient struct {
		client   *http.Client
		settings *MetaSettings
	}
)

func NewMetaHttpClient(s *MetaSettings) ports.MetaClient {
	return &MetaHttpClient{
		client:   NewHttpClient(time.Second * s.Timeout),
		settings: s,
	}
}

func NewMetaSettings(v *viper.Viper) *MetaSettings {
	return &MetaSettings{
		Retries: uint8(v.GetUint16("meta.retries")),
		Host:    v.GetString("meta.host"),
		Timeout: v.GetDuration("meta.timeout"),
	}
}

func (m *MetaHttpClient) ReadMessage(ctx context.Context,
	t *domain.Tenant,
	id string) error {
	u := m.parseUrl(fmt.Sprintf("/v16.0/%s/messages", t.AccountSettings.PhoneId))
	rmessage := boundaries.NewReadMessageRequest(id)
	bts, err := json.Marshal(rmessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(bts))
	req.Header.Add(AUTH_HEADER, fmt.Sprintf("Bearer %s", t.AccountSettings.Token))
	req.Header.Add(CONTENT_TYPE, "application/json")
	if err != nil {
		return err
	}
	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		defer res.Body.Close()
		btsr, err := io.ReadAll(res.Body)
		if err != nil {
			return nil
		}
		return errors.New(string(btsr))
	}
	return nil
}

func (m *MetaHttpClient) SendTextMessage(ctx context.Context,
	t *domain.Tenant,
	to string,
	b string) (string, error) {
	u := m.parseUrl(fmt.Sprintf("/v16.0/%s/messages", t.AccountSettings.PhoneId))
	msg := boundaries.NewSendTextMessageRequest(to, b)
	bts, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	fmt.Println(string(bts))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(bts))
	req.Header.Add(AUTH_HEADER, fmt.Sprintf("Bearer %s", t.AccountSettings.Token))
	req.Header.Add(CONTENT_TYPE, "application/json")
	if err != nil {
		return "", err
	}
	res, err := m.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	btsr, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(btsr))
	if res.StatusCode > 299 {
		return "", errors.New("META RESPONDS WITH: " + string(btsr))
	}
	msgres := &boundaries.SentMessageOutput{}
	if err = json.Unmarshal(btsr, msgres); err != nil {
		return "", err
	}
	return msgres.Messages[0].Id, nil
}

func (m *MetaHttpClient) parseUrl(path string) *url.URL {
	u, err := url.Parse(m.settings.Host)
	if err != nil {
		panic(err)
	}
	return u.JoinPath(path)
}
