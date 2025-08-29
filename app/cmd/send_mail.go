package cmd

import (
	"bytes"
	"log/slog"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/szks-repo/small-business-agents/app/pkg/mail/maildata"
	smtplib "github.com/szks-repo/small-business-agents/app/pkg/smtp"
)

var sendMailCmd = &cobra.Command{
	Use: "sendMail",
	Run: func(cmd *cobra.Command, args []string) {

		auth := smtplib.NoAuth{Auth: smtp.PlainAuth(
			"",
			"your_email@example.com",
			"your_password",
			"localhost",
		)}

		testMail := maildata.GetRandom()

		var msgBuf bytes.Buffer
		msgBuf.WriteString(makeHeader(testMail.To, testMail.From, testMail.SenderName, testMail.Subject))
		msgBuf.WriteString("\r\n")
		body := strings.ReplaceAll(testMail.Body, "\r\n", "\n")
		body = strings.ReplaceAll(body, "\n", "\r\n")
		msgBuf.WriteString(body)

		hostPort := os.Getenv("SMTP_HOST_PORT")
		slog.Info("smtp.Addr", "hostPort", hostPort)
		if err := smtp.SendMail(hostPort, auth, testMail.From, testMail.To, msgBuf.Bytes()); err != nil {
			slog.Error("Failed to SendMail", "error", err)
		}
	},
}

func makeHeader(to []string, from, senderName, subject string) string {
	addr := mail.Address{
		Name:    senderName,
		Address: from,
	}
	var buf strings.Builder
	buf.WriteString(strings.Join([]string{
		"Content-Type: text/plain; charset=utf-8",
		"Date: " + time.Now().Format(time.RFC1123Z),
		"From: " + addr.String(),
		"To: " + strings.Join(to, ", "),
		"Subject: " + subject,
	}, "\r\n"))

	return buf.String()
}

func init() {
	rootCmd.AddCommand(sendMailCmd)
}
