package simpay

import (
	"crypto"
	"fmt"
	"strings"
	"time"
)

const dateTimeFormat = "2006-01-02T15:04:05-07:00"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	var err error
	t.Time, err = time.Parse(dateTimeFormat, s)
	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(dateTimeFormat) + `"`), nil
}

type Response struct {
	Success bool                `json:"success"`
	Errors  map[string][]string `json:"errors"`
}

type PaginatedResponse struct {
	Response
	Pagination struct {
		Total       int `json:"total"`
		Count       int `json:"count"`
		PerPage     int `json:"per_page"`
		CurrentPage int `json:"current_page"`
		TotalPages  int `json:"total_pages"`
		Links       struct {
			NextPage string `json:"next_page"`
			PrevPage string `json:"prev_page"`
		} `json:"links"`
	} `json:"pagination"`
}

type SmsServiceListResponse struct {
	PaginatedResponse
	ServiceList []struct {
		Id        int    `json:"id"`
		Type      string `json:"type"`
		Status    string `json:"status"`
		Name      string `json:"name"`
		Prefix    string `json:"prefix"`
		Suffix    string `json:"suffix"`
		Adult     bool   `json:"adult"`
		CreatedAt Time   `json:"created_at"`
	} `json:"data"`
}

type SmsServiceDetailsResponse struct {
	Response
	ServiceDetails struct {
		Id          int      `json:"id"`
		Type        string   `json:"type"`
		Status      string   `json:"status"`
		Name        string   `json:"name"`
		Prefix      string   `json:"prefix"`
		Suffix      string   `json:"suffix"`
		Description string   `json:"description"`
		Adult       bool     `json:"adult"`
		Numbers     []string `json:"numbers"`
		CreatedAt   Time     `json:"created_at"`
	} `json:"data"`
}

type SmsTransactionListResponse struct {
	PaginatedResponse
	TransactionList []struct {
		Id     int    `json:"id"`
		From   int64  `json:"from"`
		Code   string `json:"code"`
		Used   bool   `json:"used"`
		SendAt Time   `json:"send_at"`
	} `json:"data"`
}

type SmsTransactionDetailsResponse struct {
	Response
	TransactionDetails struct {
		Id         int     `json:"id"`
		From       int64   `json:"from"`
		Code       string  `json:"code"`
		Used       bool    `json:"used"`
		SendNumber int     `json:"send_number"`
		Value      float64 `json:"value"`
		SendAt     Time    `json:"send_at"`
	} `json:"data"`
}

type NumberDetails struct {
	Number     int     `json:"number"`
	Value      float64 `json:"value"`
	ValueGross float64 `json:"value_gross"`
	Adult      bool    `json:"adult"`
}

type ServiceNumberListResponse struct {
	PaginatedResponse
	NumberList []NumberDetails `json:"data"`
}

type NumberDetailsResponse struct {
	Response
	NumberDetails NumberDetails `json:"data"`
}

type NumberListResponse struct {
	PaginatedResponse
	NumberList []NumberDetails `json:"data"`
}

type CodeVerifyRequest struct {
	Code   string
	Number int64
}

type CodeVerificationResponse struct {
	Success          bool `json:"success"`
	CodeVerification struct {
		Used   bool    `json:"used"`
		Code   string  `json:"code"`
		Test   bool    `json:"test"`
		From   string  `json:"from"`
		Number int     `json:"number"`
		Value  float64 `json:"value"`
		UsedAt Time    `json:"used_at"`
	} `json:"data"`
}

type DirectBillingServiceListResponse struct {
	PaginatedResponse
	ServiceList []struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Suffix    string `json:"suffix"`
		Status    string `json:"status"`
		CreatedAt Time   `json:"created_at"`
	} `json:"data"`
}

type CommissionPercent struct {
	Commission0  string `json:"commission_0"`
	Commission9  string `json:"commission_9"`
	Commission25 string `json:"commission_25"`
}

type DirectBillingServiceDetailsResponse struct {
	Response
	ServiceDetails struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Suffix string `json:"suffix"`
		Status string `json:"status"`
		Api    struct {
			Complete string `json:"complete"`
			Failure  string `json:"failure"`
		} `json:"api"`
		Providers struct {
			TMobile bool `json:"t-mobile"`
			Orange  bool `json:"orange"`
			Play    bool `json:"play"`
			Plus    bool `json:"plus"`
		} `json:"providers"`
		Commissions struct {
			TMobile CommissionPercent `json:"t-mobile"`
			Orange  CommissionPercent `json:"orange"`
			Play    CommissionPercent `json:"play"`
			Plus    CommissionPercent `json:"plus"`
		} `json:"commissions"`
		MaxValues struct {
			TMobile string `json:"t-mobile"`
			Orange  string `json:"orange"`
			Play    string `json:"play"`
			Plus    string `json:"plus"`
		} `json:"maxValues"`
		CreatedAt Time `json:"created_at"`
	} `json:"data"`
}

type CommissionValue struct {
	Net   float64 `json:"net"`
	Gross float64 `json:"gross"`
}

type CalculateCommissionResponse struct {
	Response
	CalculateCommission struct {
		Orange  CommissionValue `json:"orange"`
		Play    CommissionValue `json:"play"`
		TMobile CommissionValue `json:"t-mobile"`
		Plus    CommissionValue `json:"plus"`
	} `json:"data"`
}

type DirectBillingTransactionListResponse struct {
	PaginatedResponse
	TransactionList []struct {
		Id         string  `json:"id"`
		Status     string  `json:"status"`
		Value      float64 `json:"value"`
		ValueNetto float64 `json:"value_netto"`
		Operator   string  `json:"operator"`
		CreatedAt  Time    `json:"created_at"`
		UpdatedAt  Time    `json:"updated_at"`
	} `json:"data"`
}

type DirectBillingTransactionDetailsResponse struct {
	Response
	TransactionDetails struct {
		Id          string      `json:"id"`
		Status      string      `json:"status"`
		PhoneNumber interface{} `json:"phoneNumber"`
		Control     string      `json:"control"`
		Value       float64     `json:"value"`
		ValueNetto  float64     `json:"value_netto"`
		Operator    string      `json:"operator"`
		Notify      struct {
			IsSend     bool `json:"is_send"`
			LastSendAt Time `json:"last_send_at"`
			Count      int  `json:"count"`
		} `json:"notify"`
		CreatedAt Time `json:"created_at"`
		UpdatedAt Time `json:"updated_at"`
	} `json:"data"`
}

type GenerateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	AmountType  string  `json:"amountType"`
	Description string  `json:"description"`
	Control     string  `json:"control"`
	Returns     struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
	} `json:"returns"`
	PhoneNumber string `json:"phoneNumber"`
	Signature   string `json:"signature"`
}

func (r GenerateTransactionRequest) Sign(key string) {
	fields := []string{fmt.Sprintf("%f", r.Amount), r.AmountType, r.Description, r.Control, r.Returns.Success, r.Returns.Failure, r.PhoneNumber, key}
	r.Signature = string(crypto.SHA256.New().Sum([]byte(strings.Join(fields, "|"))))
}

func (r GenerateTransactionRequest) SignWithAmountAndControl(key string) {
	r.Signature = string(crypto.SHA256.New().Sum([]byte(strings.Join([]string{fmt.Sprintf("%f", r.Amount), r.Control, key}, "|"))))
}

type DirectBillingGenerateTransactionResponse struct {
	Response
	Data struct {
		TransactionId string `json:"transactionId"`
		RedirectUrl   string `json:"redirectUrl"`
	} `json:"data"`
}

type DirectBillingTransactionNotification struct {
	Id        string `json:"id"`
	ServiceId int    `json:"service_id"`
	Status    string `json:"status"`
	Values    struct {
		Net     float64 `json:"net"`
		Gross   float64 `json:"gross"`
		Partner float64 `json:"partner"`
	} `json:"values"`
	Returns struct {
		Complete string `json:"complete"`
		Failure  string `json:"failure"`
	} `json:"returns"`
	Control    string `json:"control"`
	NumberFrom string `json:"number_from"`
	Provider   int    `json:"provider"`
	Signature  string `json:"Signature"`
}
