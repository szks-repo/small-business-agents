package webhook

type WebhookKind int

const (
	WebhookKindContactReceived = iota + 1
	WebhookKindEmailReceived
)

var WebhookPathToKind = map[string]WebhookKind{
	"/contact/received": WebhookKindContactReceived,
	"/email/received":   WebhookKindEmailReceived,
}
