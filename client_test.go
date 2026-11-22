package simpay

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// fakeTransport answers every request with body and records the last request.
type fakeTransport struct {
	body        string
	lastRequest *http.Request
	lastBody    string
}

func (f *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	f.lastRequest = r
	if r.Body != nil {
		b, _ := io.ReadAll(r.Body)
		f.lastBody = string(b)
	}
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(f.body)),
		Request:    r,
	}, nil
}

func fakeClient(body string) (restClient, *fakeTransport) {
	transport := &fakeTransport{body: body}
	return restClient{http.Client{Transport: interceptor{apikey: "key", simKey: "password", core: transport}}}, transport
}

// Response bodies below follow the shape of real SimPay API responses.

func TestVerifyCodeSendsLowercaseFields(t *testing.T) {
	client, transport := fakeClient(`{"success":true,"data":{"used":false,"code":"ABC123","test":false,"from":"48123456789","number":7136,"value":"0.50"}}`)
	sms := Sms{restClient: client}

	response, err := sms.VerifyCode(1, "ABC123", 7136)
	if err != nil {
		t.Fatal(err)
	}

	var sent map[string]interface{}
	if err := json.Unmarshal([]byte(transport.lastBody), &sent); err != nil {
		t.Fatal(err)
	}
	if sent["code"] != "ABC123" || sent["number"] != float64(7136) {
		t.Errorf("request body = %s, want code and number in lowercase", transport.lastBody)
	}
	if !response.Success || response.CodeVerification.From != "48123456789" || response.CodeVerification.Number != 7136 || response.CodeVerification.Value != 0.5 {
		t.Errorf("unexpected response %+v", response)
	}
}

func TestNumberListReadsDecimalStringValues(t *testing.T) {
	client, _ := fakeClient(`{"success":true,"data":[{"number":7055,"value":"0.50","value_gross":0.62,"adult":false},{"number":7255,"value":2,"value_gross":2.46,"adult":false}],"pagination":{"total":2,"count":2,"per_page":15,"current_page":1,"total_pages":1,"links":{"next_page":null,"prev_page":null}}}`)
	sms := Sms{restClient: client}

	response, err := sms.GetNumbers(1, 15)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.NumberList) != 2 || response.NumberList[0].Value != 0.5 || response.NumberList[1].Value != 2 {
		t.Errorf("unexpected numbers %+v", response.NumberList)
	}
}

func TestVerifyCodeOmitsNumberWhenNotGiven(t *testing.T) {
	client, transport := fakeClient(`{"success":true,"data":{"used":true,"code":"ABC123","test":false,"from":"48123456789","number":7136,"value":0.5,"used_at":"2026-11-22T10:00:00+01:00"}}`)
	sms := Sms{restClient: client}

	response, err := sms.VerifyCode(1, "ABC123", 0)
	if err != nil {
		t.Fatal(err)
	}
	if transport.lastBody != `{"code":"ABC123"}` {
		t.Errorf("request body = %s", transport.lastBody)
	}
	if got := response.CodeVerification.UsedAt.Format("2006-01-02 15:04"); got != "2026-11-22 10:00" {
		t.Errorf("used_at = %s", got)
	}
}

func TestServiceDetailsAcceptsShortUuidId(t *testing.T) {
	client, _ := fakeClient(`{"success":true,"data":{"id":"d151e4f9","type":"ONE_TIME_CODE","status":"service_active","name":"Sklep","prefix":"SIM","suffix":"SKLEP","description":null,"adult":false,"numbers":["7136","92555"],"created_at":"2026-01-01T00:00:00+01:00"}}`)
	sms := Sms{restClient: client}

	response, err := sms.GetServiceDetails(1)
	if err != nil {
		t.Fatal(err)
	}
	if response.ServiceDetails.Id != "d151e4f9" {
		t.Errorf("id = %q", response.ServiceDetails.Id)
	}
}

func TestDirectBillingServiceListAcceptsShortUuidId(t *testing.T) {
	client, _ := fakeClient(`{"success":true,"data":[{"id":"d151e4f9","name":"Sklep","suffix":"SKLEP","status":"service_db_active","created_at":"2026-01-01T00:00:00+01:00"}],"pagination":{"total":1,"count":1,"per_page":15,"current_page":1,"total_pages":1,"links":{"next_page":null,"prev_page":null}}}`)
	d := DirectBilling{restClient: client}

	response, err := d.GetServiceList(1, 15)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.ServiceList) != 1 || response.ServiceList[0].Id != "d151e4f9" {
		t.Errorf("unexpected list %+v", response.ServiceList)
	}
}

func TestDirectBillingTransactionDetailsReadsNetValueAndNullDates(t *testing.T) {
	client, _ := fakeClient(`{"success":true,"data":{"id":"tx-1","status":"transaction_db_new","phoneNumber":null,"control":null,"value":12.3,"value_net":10,"value_created":10,"value_created_type":"NET","operator":null,"notify":{"is_send":false,"last_send_at":"1970-01-01T01:00:00+01:00","count":null},"score":null,"created_at":"2026-11-22T09:59:00+01:00","updated_at":null}}`)
	d := DirectBilling{restClient: client}

	response, err := d.GetTransactionDetails(1, "tx-1")
	if err != nil {
		t.Fatal(err)
	}
	if response.TransactionDetails.ValueNetto != 10 {
		t.Errorf("value net = %v", response.TransactionDetails.ValueNetto)
	}
	if !response.TransactionDetails.UpdatedAt.IsZero() {
		t.Errorf("updated_at = %v, want zero time", response.TransactionDetails.UpdatedAt)
	}
}
