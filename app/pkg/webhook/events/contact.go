package events

type ContactReceived struct {
	BusinessKind    string
	CompanyName     string
	CompanyNameKana string
	StaffName       string
	Content         string
}

func (p *ContactReceived) Unmarshal(data []byte) error {
	return nil
}
