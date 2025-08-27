package events

import "encoding/json"

type EmailReceived struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	CC     []string `json:"cc"`
	Header string   `json:"header"`
	Body   string   `json:"body"`
}

func (p *EmailReceived) Unmarshal(data []byte) error {
	//TBD
	var dst EmailReceived
	json.Unmarshal(data, &dst)
	p = &dst

	return nil
}
