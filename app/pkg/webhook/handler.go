package webhook

import (
	"context"
	"log/slog"

	"github.com/szks-repo/small-business-agents/app/app/pkg/types"
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
		//
	case WebhookKindEmailReceived:
		//
	default:
		slog.Info("Unknown webhook kind", "path", payload.Path)
	}

	return nil
}
