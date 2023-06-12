package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
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

func (m *MetaHttpClient) getMediaContent(ctx context.Context, url string, tk string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Add(AUTH_HEADER, fmt.Sprintf("Bearer %s", tk))
	if err != nil {
		return nil, err
	}
	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		r, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(r))
	}
	return io.ReadAll(res.Body)
}

func (m *MetaHttpClient) GetAudio(ctx context.Context,
	t *domain.Tenant,
	id string) ([]byte, error) {
	u := m.parseUrl(fmt.Sprintf("/v16.0/%s", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	req.Header.Add(AUTH_HEADER, fmt.Sprintf("Bearer %s", t.AccountSettings.Token))
	if err != nil {
		return nil, err
	}
	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	btsr, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, errors.New(string(btsr))
	}
	media := &boundaries.Media{}
	if err := json.Unmarshal(btsr, media); err != nil {
		return nil, err
	}
	return m.getMediaContent(ctx, media.URL, t.AccountSettings.Token)
}

func createAudioFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", "audio/ogg")
	return w.CreatePart(h)
}

func (m *MetaHttpClient) UploadMedia(ctx context.Context, t *domain.Tenant, stream io.Reader) (string, error) {
	u := m.parseUrl(fmt.Sprintf("/v16.0/%s/media", t.AccountSettings.PhoneId))
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	defer writer.Close()
	pw, err := createAudioFormFile(writer, "temp.ogg")
	if err != nil {
		return "", nil
	}
	if _, err := io.Copy(pw, stream); err != nil {
		return "", err
	}
	if err := writer.WriteField("type", "audio/ogg"); err != nil {
		return "", err
	}
	if err := writer.WriteField("messaging_product", "whatsapp"); err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), payload)
	if err != nil {
		return "", err
	}
	req.Header.Add(AUTH_HEADER, fmt.Sprintf("Bearer %s", t.AccountSettings.Token))
	req.Header.Set(CONTENT_TYPE, writer.FormDataContentType())
	res, err := m.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 {
		return "", errors.New(string(bts))
	}
	var mp struct {
		Id string `json:"id"`
	}
	if err := json.Unmarshal(bts, &mp); err != nil {
		return "", err
	}
	return mp.Id, nil
}

func (m *MetaHttpClient) SendAudioMessage(ctx context.Context,
	t *domain.Tenant,
	to string,
	id string) (string, error) {
	msg := boundaries.NewSendAudioMessageRequest(to, id)
	return m.sendMessage(ctx, t, msg)
}

func (m *MetaHttpClient) sendMessage(ctx context.Context, t *domain.Tenant, msg interface{}) (string, error) {
	u := m.parseUrl(fmt.Sprintf("/v16.0/%s/messages", t.AccountSettings.PhoneId))
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

func (m *MetaHttpClient) SendTextMessage(ctx context.Context,
	t *domain.Tenant,
	to string,
	b string) (string, error) {
	msg := boundaries.NewSendTextMessageRequest(to, b)
	return m.sendMessage(ctx, t, msg)
}

func (m *MetaHttpClient) parseUrl(path string) *url.URL {
	u, err := url.Parse(m.settings.Host)
	if err != nil {
		panic(err)
	}
	return u.JoinPath(path)
}
