package events

type EmailReceived struct {
	From   string
	To     string
	CC     []string
	Header string
	Body   string
}

func (p *EmailReceived) Unmarshal(data []byte) error {
	return nil
}
