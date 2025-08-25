package cmd

import (
	"io"
	"log"
	"net/http"
	"os"

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

		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using environment variables")
		}

		// endpointURL := os.Getenv("AWS_ENDPOINT_URL")
		region := os.Getenv("AWS_REGION")
		queueURL := os.Getenv("SQS_QUEUE_URL")

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
		if err != nil {
			panic(err)
		}
		sqsClient := sqs.NewFromConfig(cfg)

		mux := http.NewServeMux()
		mux.HandleFunc("POST /webhook", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
				MessageBody: aws.String(string(body)),
				QueueUrl:    &queueURL,
			})
		})
		srv := http.Server{Addr: ":3000", Handler: mux}
		srv.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(webhookReceiverCmd)
}
