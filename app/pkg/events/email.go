package events

import (
	"encoding/json"
)

type EmailReceived struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	// Header  string   `json:"header"`
}

func (p *EmailReceived) Unmarshal(data []byte) error {
	var dst EmailReceived
	json.Unmarshal(data, &dst)
	p.From = dst.From
	p.To = dst.To
	p.CC = dst.CC
	p.Subject = dst.Subject
	p.Body = dst.Body
	return nil
}
