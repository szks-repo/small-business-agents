package maildata

import (
	_ "embed"
	"math/rand/v2"

	"github.com/goccy/go-yaml"
)

type MailData struct {
	Case       string   `yaml:"case"`
	From       string   `yaml:"from"`
	SenderName string   `yaml:"senderName"`
	To         []string `yaml:"to"`
	Subject    string   `yaml:"subject"`
	Body       string   `yaml:"body"`
}

var mailData []MailData

//go:embed maildata.yaml
var mailSample []byte

func init() {
	yaml.Unmarshal(mailSample, &mailData)
}

func GetRandom() MailData {
	return mailData[rand.IntN(len(mailData))]
}
