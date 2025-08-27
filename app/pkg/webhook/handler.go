package webhook

import (
	"context"
	"fmt"

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
	switch WebhookPathToKind[payload.Path] {
	case WebhookKindContactReceived:
		return h.handleContactReceived(ctx, payload.Body)
	case WebhookKindEmailReceived:
		return h.handleEmailReceived(ctx, payload.Body)
	default:
		return fmt.Errorf("Unknown webhook kind: path=%s", payload.Path)
	}
}

func (h *webhookHandler) handleContactReceived(ctx context.Context, body []byte) error {
	var e events.ContactReceived
	if err := e.Unmarshal(body); err != nil {
		return err
	}

	return nil
}

func (h *webhookHandler) handleEmailReceived(ctx context.Context, body []byte) error {
	var e events.EmailReceived
	if err := e.Unmarshal(body); err != nil {
		return err
	}

	return nil
}
