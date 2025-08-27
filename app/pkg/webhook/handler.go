package webhook

import (
	"context"
	"log/slog"

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
		return h.handleContactReceived(ctx, payload.Body)
	default:
		slog.Info("Unknown webhook kind", "path", payload.Path)
		return nil
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
