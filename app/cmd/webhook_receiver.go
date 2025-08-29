package cmd

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/emersion/go-smtp"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/szks-repo/small-business-agents/app/pkg/events"
	smtplib "github.com/szks-repo/small-business-agents/app/pkg/smtp"
	"github.com/szks-repo/small-business-agents/app/pkg/types"
)

var webhookReceiverCmd = &cobra.Command{
	Use: "webhookReceiver",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		nctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		defer stop()

		endpoint := os.Getenv("SQS_ENDPOINT")
		region := os.Getenv("AWS_REGION")

		cfg := lo.Must(config.LoadDefaultConfig(ctx, config.WithRegion(region)))

		sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})

		queueUrl := lo.Must(sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String("webhook-event-queue"),
		}))

		mux := http.NewServeMux()
		mux.HandleFunc("POST /webhook", func(w http.ResponseWriter, r *http.Request) {
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
				MessageBody:    aws.String(string(msgBody)),
				QueueUrl:       queueUrl.QueueUrl,
				MessageGroupId: aws.String(rand.Text()),
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		srv := http.Server{Addr: os.Getenv("WEBHOOK_RECEIVER_ADDR"), Handler: mux}
		go func() {
			slog.Info("HTTP Server Started")
			if err := srv.ListenAndServe(); err != nil {
				slog.Warn("ListenAndServe", "error", err)
			}
		}()

		smtpSrv := smtp.NewServer(smtplib.NewBackend(func(r io.Reader) {
			slog.Info("DataHandler Started")
			msg, err := mail.ReadMessage(r)
			if err != nil {
				slog.Error("Failed to ReadMessage", "error", err)
				return
			}

			body, err := io.ReadAll(msg.Body)
			if err != nil {
				slog.Error("Failed to ReadAll", "error", err)
				return
			}

			res, err := http.Post(
				"http://localhost"+os.Getenv("WEBHOOK_RECEIVER_ADDR")+"/webhook",
				"application/json",
				bytes.NewReader(lo.Must(json.Marshal(events.EmailReceived{
					From:    msg.Header.Get("From"),
					To:      strings.Split(msg.Header.Get("To"), ", "),
					CC:      strings.Split(msg.Header.Get("CC"), ", "),
					Subject: msg.Header.Get("Subject"),
					Body:    string(body),
				}))),
			)
			if err != nil {
				slog.Error("Failed to Post", "error", err)
				return
			}
			defer func() {
				io.Copy(io.Discard, res.Body)
				res.Body.Close()
			}()

			slog.Info("Success Post", "statusCode", res.StatusCode)
		}))
		smtpSrv.Addr = ":2500"
		smtpSrv.Domain = "localhost"

		go func() {
			slog.Info("Starting SMTP server on :2500...")
			if err := smtpSrv.ListenAndServe(); err != nil {
				slog.Error("ListenAndServe", "error", err)
			}
		}()

		<-nctx.Done()
		slog.Info("Received shutdown signal, stopping worker")

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		smtpSrv.Shutdown(ctx)
		slog.Info("Worker stopped gracefully")
	},
}

func init() {
	rootCmd.AddCommand(webhookReceiverCmd)
}
