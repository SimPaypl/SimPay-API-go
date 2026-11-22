package simpay

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type DirectBilling struct {
	restClient
}

func NewDirectBilling(apiKey, simPassword string) DirectBilling {
	return DirectBilling{
		restClient: newRestClient(apiKey, simPassword),
	}
}

func (d DirectBilling) GetServiceList(page, limit uint) (DirectBillingServiceListResponse, error) {
	endpoint := fmt.Sprintf("/directbilling?page=%v&limit=%v", page, limit)
	response, err := d.restClient.sendGetRequest(endpoint)
	if err != nil {
		return DirectBillingServiceListResponse{}, err
	}
	var serviceList DirectBillingServiceListResponse
	return serviceList, json.Unmarshal(response, &serviceList)
}
func (d DirectBilling) GetServiceDetails(serviceId uint) (DirectBillingServiceDetailsResponse, error) {
	endpoint := fmt.Sprintf("/directbilling/%v", serviceId)
	response, err := d.restClient.sendGetRequest(endpoint)
	if err != nil {
		return DirectBillingServiceDetailsResponse{}, err
	}
	var serviceDetails DirectBillingServiceDetailsResponse
	return serviceDetails, json.Unmarshal(response, &serviceDetails)
}
func (d DirectBilling) CalculateCommission(serviceId uint, amount float32) (CalculateCommissionResponse, error) {
	endpoint := fmt.Sprintf("/directbilling/%v/calculate?amount=%f", serviceId, amount)
	response, err := d.restClient.sendGetRequest(endpoint)
	if err != nil {
		return CalculateCommissionResponse{}, err
	}
	var commissionResponse CalculateCommissionResponse
	return commissionResponse, json.Unmarshal(response, &commissionResponse)
}

func (d DirectBilling) GetTransactionList(serviceId, page, limit uint) (DirectBillingTransactionListResponse, error) {
	endpoint := fmt.Sprintf("/directbilling/%v/transactions?page=%v&limit=%v", serviceId, page, limit)
	response, err := d.restClient.sendGetRequest(endpoint)
	if err != nil {
		return DirectBillingTransactionListResponse{}, err
	}
	var transactionList DirectBillingTransactionListResponse
	return transactionList, json.Unmarshal(response, &transactionList)
}
func (d DirectBilling) GetTransactionDetails(serviceId uint, transactionId string) (DirectBillingTransactionDetailsResponse, error) {
	endpoint := fmt.Sprintf("/directbilling/%v/transactions/%v", serviceId, transactionId)
	response, err := d.restClient.sendGetRequest(endpoint)
	if err != nil {
		return DirectBillingTransactionDetailsResponse{}, err
	}
	var transactionDetails DirectBillingTransactionDetailsResponse
	return transactionDetails, json.Unmarshal(response, &transactionDetails)
}
func (d DirectBilling) GenerateTransaction(serviceId uint, request GenerateTransactionRequest) (DirectBillingGenerateTransactionResponse, error) {
	endpoint := fmt.Sprintf("/directbilling/%d/transactions", serviceId)
	response, err := d.restClient.sendPostRequest(endpoint, request)
	if err != nil {
		return DirectBillingGenerateTransactionResponse{}, err
	}
	var transactionResponse DirectBillingGenerateTransactionResponse
	return transactionResponse, json.Unmarshal(response, &transactionResponse)
}

// CheckSignature reports whether transactionJson, the body of a direct billing
// notification, carries a valid signature for key.
//
// SimPay signs the values of id, service_id, status, values.net, values.gross,
// values.partner, returns.complete, returns.failure, control, number_from and
// provider, in that order, skipping fields that are absent, followed by key,
// joined with "|" and hashed with SHA-256. Numbers are taken exactly as they
// appear in the JSON body.
func CheckSignature(key, transactionJson string) bool {
	decoder := json.NewDecoder(strings.NewReader(transactionJson))
	decoder.UseNumber()

	var n map[string]interface{}
	if err := decoder.Decode(&n); err != nil {
		return false
	}

	signature, ok := n["signature"].(string)
	if !ok || signature == "" {
		return false
	}

	values, _ := n["values"].(map[string]interface{})
	returns, _ := n["returns"].(map[string]interface{})

	var fields []string
	for _, value := range []interface{}{
		n["id"],
		n["service_id"],
		n["status"],
		values["net"],
		values["gross"],
		values["partner"],
		returns["complete"],
		returns["failure"],
		n["control"],
		n["number_from"],
		n["provider"],
	} {
		if value != nil {
			fields = append(fields, fmt.Sprint(value))
		}
	}
	fields = append(fields, key)

	sum := sha256.Sum256([]byte(strings.Join(fields, "|")))
	expected := hex.EncodeToString(sum[:])

	return subtle.ConstantTimeCompare([]byte(expected), []byte(strings.ToLower(signature))) == 1
}
