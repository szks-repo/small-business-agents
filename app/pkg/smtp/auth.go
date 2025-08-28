package smtplib

import "net/smtp"

type NoAuth struct{ smtp.Auth }

func (na NoAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	srv := *server
	srv.TLS = true
	return na.Auth.Start(&srv)
}
