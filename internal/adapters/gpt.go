package adapters

import (
	"context"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

const (
	MAIN_DESCRIPTION = `Considere "SEMPRE" essa legenda para gerar seus textos. 
	#[USER] Indica que o texto a seguir vem do usuario final da aplicacao.
	#[APPLICATION] Indica que o texto é uma resposta ja gerada por voce, gere seus textos "SEMPRE" com esse prefixo: #[APPLICATION]`
)

type (
	GptClient struct {
		token string
	}
)

func NewGptClient(v *viper.Viper) ports.NlpClient {
	tk := v.GetString("gpt.token")
	return &GptClient{
		token: tk,
	}
}

func (gc *GptClient) UnrollConversation(ctx context.Context, msgs []domain.Msg) (*domain.Msg, error) {
	p := MAIN_DESCRIPTION + domain.CompileMessages(msgs)
	req := openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		Prompt:      p,
		Temperature: 1,
		MaxTokens:   512,
	}
	cli := openai.NewClient(gc.token)
	res, err := cli.CreateCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	txt := res.Choices[0].Text
	txt = strings.ReplaceAll(txt, p, "")
	txt = strings.ReplaceAll(txt, string(domain.APPLICATION), "")
	txt = strings.ReplaceAll(txt, string(domain.USER), "")
	txt = strings.ReplaceAll(txt, "#", "")
	return &domain.Msg{
		Role:    domain.APPLICATION,
		Content: txt,
		Status:  "generated",
	}, nil
}
