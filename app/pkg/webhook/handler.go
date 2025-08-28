package webhook

import (
	"context"

	"github.com/szks-repo/small-business-agents/app/app/pkg/types"
	"github.com/szks-repo/small-business-agents/app/app/pkg/webhook/events"
)

type WebhookHandler interface {
	Handle(ctx context.Context, payload *types.WebhookPayload) error
}

type webhookHandler struct {
}

func NewHandler() WebhookHandler {
	return &webhookHandler{}
}

func (h *webhookHandler) Handle(ctx context.Context, payload *types.WebhookPayload) error {
	var event events.EmailReceived
	if err := event.Unmarshal(payload.Body); err != nil {
		return err
	}

	//TODO

	return nil
}
