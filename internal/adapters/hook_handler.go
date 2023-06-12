package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pm2/internal/domain"
	"pm2/internal/ports"
	"pm2/internal/ports/boundaries"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	WhatsAppHookHandler struct {
		conversationService       ports.ConversationUseCase
		incomingMessageRepository ports.Repository[domain.IncomingMessage]
	}
)

func NewWhatsAppHookHandler(conversationService ports.ConversationUseCase,
	incomingMessageRepository ports.Repository[domain.IncomingMessage]) *WhatsAppHookHandler {
	return &WhatsAppHookHandler{
		conversationService:       conversationService,
		incomingMessageRepository: incomingMessageRepository,
	}
}

func (wh *WhatsAppHookHandler) CheckHook(ctx *gin.Context) {
	c := ctx.Query("hub.challenge")
	cc, _ := strconv.ParseInt(c, 10, 0)
	ctx.JSON(http.StatusOK, cc)
}

func (wh *WhatsAppHookHandler) IncomingMessage(ctx *gin.Context) {
	b, _ := io.ReadAll(ctx.Request.Body)
	m := boundaries.IncomingMessageInput{}
	s := string(b)
	fmt.Println(s)
	if err := json.Unmarshal(b, &m); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if len(m.Entry[0].Changes[0].Value.Messages) == 0 {
		ctx.Status(http.StatusNoContent)
		return
	}
	c := m.Entry[0].Changes[0]
	ms := c.Value.Messages[0]
	rm := wh.incomingMessageRepository.GetFirst(ctx, ports.GetById(ms.Id))
	if rm != nil {
		ctx.JSON(http.StatusOK, m)
		return
	}
	wh.incomingMessageRepository.Insert(ctx, domain.NewIncomingMessage(ms.Id, m))
	wh.conversationService.UnrollConversation(ctx, &m)
	ctx.JSON(http.StatusOK, m)
}
