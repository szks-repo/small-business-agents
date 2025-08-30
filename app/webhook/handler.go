package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/samber/lo"

	"github.com/szks-repo/small-business-agents/app/pkg/events"
	mailllib "github.com/szks-repo/small-business-agents/app/pkg/mail"
	"github.com/szks-repo/small-business-agents/app/pkg/types"
)

type WebhookHandler interface {
	Handle(ctx context.Context, payload *types.WebhookPayload) error
}

type webhookHandler struct {
	fastApiUrl string
}

func NewHandler() WebhookHandler {
	return &webhookHandler{
		fastApiUrl: "http://localhost:8000",
	}
}

type Request struct {
	From       string   `json:"from_address"`
	SenderName string   `json:"sender_name"`
	To         []string `json:"to"`
	CC         []string `json:"cc"`
	// Header  mail.Header   `json:"header"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
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

	jsonBody := lo.Must(json.Marshal(&Request{
		From:       addr.Address,
		SenderName: addr.Name,
		To:         event.To,
		CC:         event.To,
		// Header:     event.Header,
		Subject: event.Subject,
		Body:    event.Body,
	}))

	fmt.Println(string(jsonBody))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.fastApiUrl+"/inbox", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	slog.Info("Http Request", "statusCode", res.StatusCode, "body", body)

	return nil
}
