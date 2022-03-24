package simpay

import (
	"crypto"
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

func CheckSignature(key, transactionJson string) bool {
	var n DirectBillingTransactionNotification
	err := json.Unmarshal([]byte(transactionJson), &n)
	if err != nil {
		return false
	}

	fields := []string{
		fmt.Sprintf("%v", n.Id),
		n.Status,
		fmt.Sprintf("%f", n.Values.Net),
		fmt.Sprintf("%f", n.Values.Gross),
		fmt.Sprintf("%f", n.Values.Partner),
		n.Returns.Complete,
		n.Returns.Failure,
		n.NumberFrom,
		fmt.Sprintf("%v", n.Provider),
		n.Signature,
		key,
	}
	return string(crypto.SHA256.New().Sum([]byte(strings.Join(fields, "|")))) == n.Signature

}
