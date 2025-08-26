package cmd

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var webhookReceiverCmd = &cobra.Command{
	Use: "webhookReceiver",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		_, file, _, _ := runtime.Caller(0)
		envPath := filepath.Join(file, "../../../.env")
		if err := godotenv.Load(envPath); err != nil {
			slog.Info("No .env file found, using environment variables")
		}

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

			msgBody, err := json.Marshal(map[string]any{
				"path": r.URL.Path,
				"body": string(body),
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

		srv := http.Server{Addr: ":3000", Handler: mux}
		slog.Info("HTTP Server Started")
		srv.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(webhookReceiverCmd)
}
