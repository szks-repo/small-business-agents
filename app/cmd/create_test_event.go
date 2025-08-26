package cmd

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var createTestEventPath string
var createTestEventMessage string

var createTestEventCmd = &cobra.Command{
	Use: "createTestEvent",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("createTestEvent start")

		ctx := cmd.Context()

		payload := lo.Must(json.Marshal(map[string]any{
			"ts":      time.Now(),
			"message": createTestEventMessage,
		}))

		webhookUrl := lo.Must(url.JoinPath("http://localhost"+os.Getenv("WEBHOOK_RECEIVER_ADDR"), "webhook", createTestEventPath))
		req, err := http.NewRequestWithContext(ctx, "POST", webhookUrl, bytes.NewReader(payload))
		if err != nil {
			panic(err)
		}
		client := http.DefaultClient
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			slog.Error("Failed to send webhook", "webhookUrl", webhookUrl, "statusCode", res.StatusCode)
			return
		}

		slog.Info("Success to send webhook", "statusCode", res.StatusCode, "payload", string(payload))
	},
}

func init() {
	createTestEventCmd.Flags().StringVar(&createTestEventPath, "p", "/contact", "")
	createTestEventCmd.Flags().StringVar(&createTestEventMessage, "m", "my message", "")
	rootCmd.AddCommand(createTestEventCmd)
}
