package adapters

import (
	"context"
	"errors"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"strings"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

const (
	MAIN_DESCRIPTION = `Seu nome é Bob, um chatbot da empresa Europ Assistance Brasil, que realiza atendimentos para contratação de seguros viagem servindo de apoio para o e-comerce: 'https://eaviagem.com.br'.

	Considere "SEMPRE" essa LEGENDA para gerar seus textos!
	#[USER] Indica que o texto a seguir vem do usuario final da aplicacao.
	#[APPLICATION] Indica que o texto é uma resposta ja gerada por voce, gere seus textos "SEMPRE" com esse prefixo: #[APPLICATION]`
)

type (
	GptClient struct {
		tenantRepository ports.Repository[domain.Tenant]
	}
)

func NewGptClient(tenantRepository ports.Repository[domain.Tenant]) ports.NlpClient {
	return &GptClient{
		tenantRepository: tenantRepository,
	}
}

func (gc *GptClient) UnrollConversation(ctx context.Context, tenantId uuid.UUID, msgs []domain.Msg) (*domain.Msg, error) {
	p := MAIN_DESCRIPTION + domain.CompileMessages(msgs)
	req := openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		Prompt:      p,
		Temperature: 1,
		MaxTokens:   512,
	}
	tenant := gc.tenantRepository.GetFirst(ctx, ports.GetById(tenantId))
	if tenant == nil {
		return nil, errors.New(domain.TENANT_NOT_FOUND)
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
	return domain.NewMessage("",
		domain.APPLICATION,
		txt,
		domain.GENERATED,
		nil), nil
}
