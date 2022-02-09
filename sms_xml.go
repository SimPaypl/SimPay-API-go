package simpay

import (
	"crypto"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const charset = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

var codes map[string]float64

func init() {
	codes = make(map[string]float64)
	codes["7055"] = 0.25
	codes["7136"] = 0.5
	codes["7255"] = 1.0
	codes["7355"] = 1.5
	codes["7455"] = 2.0
	codes["7555"] = 2.5
	codes["7636"] = 3.0
	codes["77464"] = 3.5
	codes["78464"] = 4.0
	codes["7936"] = 4.5
	codes["91055"] = 5.0
	codes["91155"] = 5.5
	codes["91455"] = 7.0
	codes["91664"] = 8.0
	codes["91955"] = 9.5
	codes["92055"] = 10.0
	codes["92555"] = 12.5
}

type SmsXml struct {
	hashingKey string
}

func NewSmsXml(hashingKey string) SmsXml {
	return SmsXml{hashingKey: hashingKey}
}

func (SmsXml) GenerateCode() string {
	var sb strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		sb.WriteString(string(charset[rand.Intn(len(charset))]))
	}
	return sb.String()
}

func (SmsXml) GenerateXml(code string) string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?><sms-response>" + code + "<sms-text></sms-text></sms-response>"
}

func (s SmsXml) CheckParameters(m map[string]interface{}) bool {
	params := []string{"send_number", "sms_text", "sms_from", "sms_id", "sign"}

	for _, param := range params {
		if !contains(param, m) {
			return false
		}
	}

	return m["sign"] == sign(m, s.hashingKey)
}

func (SmsXml) GetNumberValue(number string) float64 {
	return codes[number]
}

func contains(v string, m map[string]interface{}) bool {
	for s := range m {
		if v == s {
			return true
		}
	}
	return false
}

func sign(m map[string]interface{}, hashingKey string) string {
	values := fmt.Sprintf("%v%v%v%v%v%v", m["sms_id"], m["sms_text"], m["sms_from"], m["send_number"], m["send_time"], hashingKey)
	return string(crypto.SHA256.New().Sum([]byte(values)))
}
