package mailllib

import (
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

	addr, err := mail.ParseAddress(decodedFrom)
	if err != nil {
		return nil, err
	}

	return addr, nil
}
