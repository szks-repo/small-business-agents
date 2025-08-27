package webhook

type WebhookKind int

const (
	WebhookKindContactReceived = iota + 1
	WebhookKindEmailReceived
)

var WebhookPathToKind = map[string]WebhookKind{
	"/webhook/contact/received": WebhookKindContactReceived,
	"/webhook/email/received":   WebhookKindEmailReceived,
}
