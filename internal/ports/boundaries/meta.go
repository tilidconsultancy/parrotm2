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
		Type string `json:"type"`
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
