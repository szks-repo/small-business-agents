package webhook

type WebhookKind int

const (
	WebhookKindContactReceived = iota + 1
	WebhookKindEmailReceived
)

var WebhookPath = map[WebhookKind]string{
	WebhookKindContactReceived: "/contact/ingress",
	WebhookKindEmailReceived:   "/mail/ingress",
}
