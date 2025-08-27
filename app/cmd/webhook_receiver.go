package cmd

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/spf13/cobra"
	"github.com/szks-repo/small-business-agents/app/app/pkg/types"
)

var webhookReceiverCmd = &cobra.Command{
	Use: "webhookReceiver",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		endpoint := os.Getenv("SQS_ENDPOINT")
		region := os.Getenv("AWS_REGION")

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			panic(err)
		}

		sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})

		queueUrl, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String("webhook-event-queue"),
		})
		if err != nil {
			panic(err)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("POST /webhook/", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			msgBody, err := json.Marshal(&types.WebhookPayload{
				Path: r.URL.Path,
				Body: body,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
				MessageBody: aws.String(string(msgBody)),
				QueueUrl:    queueUrl.QueueUrl,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		srv := http.Server{Addr: os.Getenv("WEBHOOK_RECEIVER_ADDR"), Handler: mux}
		slog.Info("HTTP Server Started")
		if err := srv.ListenAndServe(); err != nil {
			slog.Warn("ListenAndServe", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webhookReceiverCmd)
}
