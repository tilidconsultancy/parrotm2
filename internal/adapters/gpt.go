package adapters

import (
	"context"
	"errors"
	"fmt"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	MAIN_DESCRIPTION = `Considere "SEMPRE" essa LEGENDA para gerar seus textos!
	#[USER] Indica que o texto a seguir vem do usuario final da aplicacao.
	#[APPLICATION] Indica que o texto é uma resposta ja gerada por voce, gere seus textos "SEMPRE" com esse prefixo: #[APPLICATION]`
)

type (
	GptClient struct {
	}
)

func NewGptClient(tenantRepository ports.Repository[domain.Tenant]) ports.NlpClient {
	return &GptClient{}
}

func (gc *GptClient) UnrollConversation(ctx context.Context, tenant *domain.Tenant, msgs []domain.Msg) (*domain.Msg, error) {
	if tenant == nil {
		return nil, errors.New(domain.TENANT_NOT_FOUND)
	}
	p := fmt.Sprintf("%s\n%s\n%s\n",
		tenant.AccountSettings.MainContext,
		MAIN_DESCRIPTION,
		domain.CompileMessages(msgs))

	req := openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		Prompt:      p,
		Temperature: 1,
		MaxTokens:   512,
	}
	cli := openai.NewClient(tenant.AccountSettings.NlpToken)
	res, err := cli.CreateCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	txt := res.Choices[0].Text
	txt = strings.ReplaceAll(txt, p, "")
	txt = strings.ReplaceAll(txt, string(domain.APPLICATION), "")
	txt = strings.ReplaceAll(txt, string(domain.USER), "")
	txt = strings.ReplaceAll(txt, "#", "")
	txt = strings.ReplaceAll(txt, "\n", "")
	return domain.NewMessage("",
		domain.APPLICATION,
		txt,
		domain.GENERATED,
		domain.TEXT,
		"",
		nil), nil
}
