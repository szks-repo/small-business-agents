package webhook

import (
	"context"

	"github.com/szks-repo/small-business-agents/app/pkg/events"
	mailllib "github.com/szks-repo/small-business-agents/app/pkg/mail"
	"github.com/szks-repo/small-business-agents/app/pkg/types"
)

type WebhookHandler interface {
	Handle(ctx context.Context, payload *types.WebhookPayload) error
}

type webhookHandler struct{}

func NewHandler() WebhookHandler {
	return &webhookHandler{}
}

func (h *webhookHandler) Handle(ctx context.Context, payload *types.WebhookPayload) error {
	var event events.EmailReceived
	if err := event.Unmarshal(payload.Body); err != nil {
		return err
	}

	addr, err := mailllib.ParseFROM(event.From)
	if err != nil {
		return err
	}

	// todo Invoke Agent1
	_ = addr

	return nil
}
