package mailllib

import (
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFROM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from    string
		want    *mail.Address
		wantErr bool
	}{
		{
			from: "info@example.com",
			want: &mail.Address{
				Address: "info@example.com",
				Name:    "",
			},
			wantErr: false,
		},
		{
			from: "Tom <info@example.com>",
			want: &mail.Address{
				Address: "info@example.com",
				Name:    "Tom",
			},
			wantErr: false,
		},
		{
			from: "=?utf-8?q?=E5=B9=B8=E3=81=9B=E6=A0=AA=E5=BC=8F=E4=BC=9A=E7=A4=BE?= <happy-corp@example.com>",
			want: &mail.Address{
				Address: "happy-corp@example.com",
				Name:    "幸せ株式会社",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := ParseFROM(tt.from)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			} else {
				if tt.wantErr {
					t.Error("error is nil")
				}
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
