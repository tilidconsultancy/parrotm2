package adapters

import (
	"context"
	"errors"
	"fmt"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"regexp"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	MAIN_DESCRIPTION = `Considere "SEMPRE" essa LEGENDA para gerar seus textos!
	#[USER] Indica que o texto a seguir vem do usuario final da aplicacao.
	#[APPLICATION] Indica que o texto é uma resposta ja gerada por voce, gere seus textos "SEMPRE" com esse prefixo: #[APPLICATION]`

	MAIN_LABEL_DESCRIPTION = `"%s" ? Dado o texto acima responda se "%s" ?, é importante que sua resposta seja apenas com o valor int32 em escala de 0 a 100.`
)

type (
	GptClient struct {
	}
)

func NewGptClient(tenantRepository ports.Repository[domain.Tenant]) ports.NlpClient {
	return &GptClient{}
}

func maprole(r domain.MsgRole) string {
	switch r {
	case domain.APPLICATION:
		return openai.ChatMessageRoleAssistant
	case domain.SYSTEM:
		return openai.ChatMessageRoleSystem
	default:
		return openai.ChatMessageRoleUser
	}
}

func unrollWithChatCompletition(ctx context.Context, tenant *domain.Tenant, msgs []domain.Msg) (*domain.Msg, error) {
	reqmsgs := []openai.ChatCompletionMessage{}
	if tenant.AccountSettings.MainContext != "" {
		reqmsgs = append(reqmsgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: tenant.AccountSettings.MainContext,
		})
	}
	for _, m := range msgs {
		reqmsgs = append(reqmsgs, openai.ChatCompletionMessage{
			Role:    maprole(m.Role),
			Content: m.Content,
		})
	}
	cli := openai.NewClient(tenant.AccountSettings.NlpToken)
	res, err := cli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Messages: reqmsgs,
		Model:    tenant.AccountSettings.Model,
	})
	if err != nil {
		return nil, err
	}
	return domain.NewMessage("",
		domain.APPLICATION,
		res.Choices[0].Message.Content,
		domain.GENERATED,
		domain.TEXT,
		"",
		nil), nil
}

func (gc *GptClient) UnrollConversation(ctx context.Context, tenant *domain.Tenant, msgs []domain.Msg) (*domain.Msg, error) {
	if tenant == nil {
		return nil, errors.New(domain.TENANT_NOT_FOUND)
	}
	if tenant.AccountSettings.ChatCompletition {
		return unrollWithChatCompletition(ctx, tenant, msgs)
	}
	p := fmt.Sprintf("%s\n%s\n%s\n",
		tenant.AccountSettings.MainContext,
		MAIN_DESCRIPTION,
		domain.CompileMessages(msgs))
	req := openai.CompletionRequest{
		Model:       tenant.AccountSettings.Model,
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

func (gc *GptClient) EvaluateLabelMeaning(
	ctx context.Context,
	tenant *domain.Tenant,
	msgs []domain.Msg,
	labels []domain.LabelMeaning) ([]domain.PercentageLabelMeaning, error) {
	r := []domain.PercentageLabelMeaning{}
	p := domain.CompileMessagesByRole(msgs, domain.USER)
	cli := openai.NewClient(tenant.AccountSettings.NlpToken)
	for _, l := range labels {
		res, err := cli.CreateCompletion(ctx, openai.CompletionRequest{
			Model:  openai.GPT3TextDavinci003,
			Prompt: fmt.Sprintf(MAIN_LABEL_DESCRIPTION, p, l.Description),
		})
		if err != nil {
			return nil, err
		}
		rg, _ := regexp.Compile("[^0-9]+")
		txt := rg.ReplaceAllString(res.Choices[0].Text, "")
		s, err := strconv.ParseInt(txt, 10, 8)
		if err != nil {
			s = 0
		}
		r = append(r, domain.PercentageLabelMeaning{
			LabelMeaning: l,
			Percentage:   uint8(s),
		})
	}
	return r, nil
}
