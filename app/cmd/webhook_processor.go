package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var webhookProcessorCmd = &cobra.Command{
	Use: "webhookProcessor",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		_, file, _, _ := runtime.Caller(0)
		envPath := filepath.Join(file, "../../../.env")
		if err := godotenv.Load(envPath); err != nil {
			slog.Info("No .env file found, using environment variables")
		}

		endpoint := os.Getenv("SQS_ENDPOINT")
		region := os.Getenv("AWS_REGION")

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
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

		slog.Info("Worker Started")
		for {
			received, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:          queueUrl.QueueUrl,
				VisibilityTimeout: 300,
				WaitTimeSeconds:   20,
			})
			if err != nil {
				slog.Error("Failed to ReceiveMessage", "error", err)
				return
			}

			for _, msg := range received.Messages {
				slog.Info("Receive message", "messageId", msg.MessageId)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(webhookProcessorCmd)
}
