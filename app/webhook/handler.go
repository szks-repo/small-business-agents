package webhook

import (
	"context"
	"fmt"

	"github.com/szks-repo/small-business-agents/app/pkg/events"
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

	//TODO
	fmt.Println("Event===>", event)

	return nil
}
