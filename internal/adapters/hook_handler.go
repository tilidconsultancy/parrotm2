package adapters

import (
	"encoding/json"
	"io"
	"net/http"
	"pm2/internal/ports"
	"pm2/internal/ports/boundaries"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	WhatsAppHookHandler struct {
		conversationService ports.ConversationUseCase
	}
)

func NewWhatsAppHookHandler(conversationService ports.ConversationUseCase) *WhatsAppHookHandler {
	return &WhatsAppHookHandler{
		conversationService: conversationService,
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
	if err := json.Unmarshal(b, &m); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if len(m.Entry[0].Changes[0].Value.Messages) == 0 {
		ctx.Status(http.StatusNoContent)
		return
	}
	wh.conversationService.UnrollConversation(ctx, &m)
	ctx.JSON(http.StatusOK, m)
}
