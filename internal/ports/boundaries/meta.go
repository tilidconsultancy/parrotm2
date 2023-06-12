package boundaries

const (
	WPP = "whatsapp"
)

type (
	SendTextMessageInput struct {
		MessageProduct string `json:"messaging_product"`
		To             string `json:"to"`
		Type           string `json:"type"`
		Text           struct {
			PreviewUrl bool   `json:"preview_url"`
			Body       string `json:"body"`
		} `json:"text"`
	}
	SendAudioMessageInput struct {
		MessageProduct string `json:"messaging_product"`
		To             string `json:"to"`
		Type           string `json:"type"`
		Audio          struct {
			Id string `json:"id"`
		} `json:"audio"`
	}
	SentMessageOutput struct {
		MessageProduct string `json:"messaging_product"`
		Messages       []struct {
			Id string `json:"id"`
		} `json:"messages"`
	}

	ReadMessageInput struct {
		MessageProduct string `json:"messaging_product"`
		Status         string `json:"status"`
		Id             string `json:"message_id"`
	}

	IncomingMessageInput struct {
		Object string `json:"object"`
		Entry  []struct {
			Id      string `json:"id"`
			Changes []struct {
				Value struct {
					MessagingProduct string `json:"messaging_product"`
					Metadata         struct {
						DisplayPhoneNumber string `json:"display_phone_number"`
						PhoneNumberId      string `json:"phone_number_id"`
					} `json:"metadata"`
					Contacts []struct {
						Profile struct {
							Name string `json:"name"`
						} `json:"profile"`
						WaId string `json:"wa_id"`
					} `json:"contacts"`
					Messages []Message `json:"messages"`
				} `json:"value"`
				Field string `json:"field"`
			} `json:"changes"`
		} `json:"entry"`
	}
	Message struct {
		From      string `json:"from"`
		Id        string `json:"id"`
		Timestamp string `json:"timestamp"`
		Text      struct {
			Body string `json:"body"`
		} `json:"text"`
		Type  string `json:"type"`
		Audio struct {
			MimeType string `json:"mime_type"`
			Sha256   string `json:"sha256"`
			Id       string `json:"id"`
			Voice    bool   `json:"voice"`
		} `json:"audio"`
	}

	Media struct {
		URL              string `json:"url"`
		MimeType         string `json:"mime_type"`
		Sha256           string `json:"sha256"`
		FileSize         int    `json:"file_size"`
		ID               string `json:"id"`
		MessagingProduct string `json:"messaging_product"`
	}
)

func NewReadMessageRequest(id string) *ReadMessageInput {
	return &ReadMessageInput{
		MessageProduct: WPP,
		Status:         "read",
		Id:             id,
	}
}

func NewSendTextMessageRequest(to string, b string) *SendTextMessageInput {
	return &SendTextMessageInput{
		MessageProduct: WPP,
		To:             to,
		Type:           "text",
		Text: struct {
			PreviewUrl bool   "json:\"preview_url\""
			Body       string "json:\"body\""
		}{
			PreviewUrl: false,
			Body:       b,
		},
	}
}

func NewSendAudioMessageRequest(to string, id string) *SendAudioMessageInput {
	return &SendAudioMessageInput{
		MessageProduct: WPP,
		To:             to,
		Type:           "audio",
		Audio: struct {
			Id string "json:\"id\""
		}{
			Id: id,
		},
	}
}
