package smtplib

import (
	"io"
	"log/slog"

	"github.com/emersion/go-smtp"
)

type DataHandler func(io.Reader)

type Backend struct {
	handler DataHandler
}

func NewBackend(h DataHandler) *Backend {
	return &Backend{handler: h}
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	slog.Info("NewSession called")
	return &Session{handler: b.handler}, nil
}

type Session struct {
	handler DataHandler
}

func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	slog.Info("Received MAIL command", "from", from)
	return nil
}

func (s *Session) Rcpt(to string, opt *smtp.RcptOptions) error {
	slog.Info("Received RCPT command", "to", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	slog.Info("Data called")
	if s.handler != nil {
		s.handler(r)
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
