package events

import (
	"encoding/json"
	"net/mail"
	"strings"
)

type EmailReceived struct {
	From    string      `json:"from"`
	To      []string    `json:"to"`
	CC      []string    `json:"cc"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
	Header  mail.Header `json:"header"`
}

func (p *EmailReceived) Unmarshal(data []byte) error {
	var dst EmailReceived
	json.Unmarshal(data, &dst)
	p.From = dst.From
	p.To = dst.To
	p.CC = dst.CC
	p.Subject = dst.Subject
	p.Body = strings.ReplaceAll(dst.Body, "\r\n", "\n")
	p.Header = dst.Header
	return nil
}
