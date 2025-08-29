package mailllib

import (
	"fmt"
	"mime"
	"net/mail"
)


func ParseFROM(from string) (*mail.Address, error) {
	// RFC 2047形式でエンコードされた文字列をデコード
	dec := new(mime.WordDecoder)
	decodedFrom, err := dec.DecodeHeader(from)
	if err != nil {
		return nil, err
	}

	// デコードした文字列をパースして、送信者名とメールアドレスを取得
	addr, err := mail.ParseAddress(decodedFrom)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Name: %s\n", addr.Name)
	fmt.Printf("Address: %s\n", addr.Address)
	return addr, nil
}
