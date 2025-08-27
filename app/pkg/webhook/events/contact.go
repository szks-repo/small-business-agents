package events

import "encoding/json"

type ContactReceived struct {
	Name            string `json:"name"`
	NameKana        string `json:"name_kana"`
	Email           string `json:"email"`
	Tel             string `json:"tel"`
	BusinessKind    string `json:"business_kind"`
	CompanyName     string `json:"company_name"`
	CompanyNameKana string `json:"company_name_kana"`
	Content         string `json:"content"`
}

func (p *ContactReceived) Unmarshal(data []byte) error {
	//TBD
	var dst ContactReceived
	json.Unmarshal(data, &dst)
	p = &dst

	return nil
}
