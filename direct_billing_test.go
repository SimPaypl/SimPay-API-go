package simpay

import (
	"strings"
	"testing"
)

// Notifications below were signed with the algorithm SimPay uses when sending
// direct billing IPNs, reproduced in PHP, with the key "secret-key".
const (
	fullNotification    = `{"id":"a1b2c3","service_id":"f1e2d3c4","status":"transaction_db_payed","values":{"net":10,"gross":12.3,"partner":7.45},"returns":{"complete":"https:\/\/shop.example\/ok","failure":"https:\/\/shop.example\/fail"},"control":"order-15","number_from":"48123456789","provider":1,"signature":"a9623d6d0a3186417f36e365ddfe8346c92025009b8cf8b1a7e3f912ec3b43ad"}`
	minimalNotification = `{"id":"d4e5f6","service_id":42,"status":"transaction_db_confirmed","values":{"net":0.5,"gross":0.62,"partner":0.25},"signature":"38ccd5c05e5a8bb5cb9d27e57a531fffbdca16500ea68d64df74ed08c4b63c29"}`
)

func TestCheckSignature(t *testing.T) {
	tests := []struct {
		name string
		key  string
		body string
		want bool
	}{
		{"all fields", "secret-key", fullNotification, true},
		{"optional fields absent, numeric service id", "secret-key", minimalNotification, true},
		{"wrong key", "other-key", fullNotification, false},
		{"tampered amount", "secret-key", strings.Replace(fullNotification, `"gross":12.3`, `"gross":123`, 1), false},
		{"tampered control", "secret-key", strings.Replace(fullNotification, `order-15`, `order-16`, 1), false},
		{"missing signature", "secret-key", `{"id":"d4e5f6","service_id":42,"status":"transaction_db_confirmed","values":{"net":0.5,"gross":0.62,"partner":0.25}}`, false},
		{"invalid json", "secret-key", `{`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSignature(tt.key, tt.body); got != tt.want {
				t.Errorf("CheckSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
