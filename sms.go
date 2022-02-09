package simpay

import (
	"encoding/json"
	"fmt"
)

type Sms struct {
	restClient restClient
}

func NewSms(apiKey, simPassword string) *Sms {
	return &Sms{
		restClient: newRestClient(apiKey, simPassword),
	}
}

func (s Sms) GetServiceList(page, limit uint) (SmsServiceListResponse, error) {
	endpoint := fmt.Sprintf("/sms?page=%v&limit=%v", page, limit)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return SmsServiceListResponse{}, err
	}
	var serviceList SmsServiceListResponse
	return serviceList, json.Unmarshal(response, &serviceList)
}

func (s Sms) GetServiceDetails(serviceId uint) (SmsServiceDetailsResponse, error) {
	endpoint := fmt.Sprintf("/sms/%v", serviceId)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return SmsServiceDetailsResponse{}, err
	}
	var serviceDetails SmsServiceDetailsResponse
	return serviceDetails, json.Unmarshal(response, &serviceDetails)
}

func (s Sms) GetTransactions(serviceId, page, limit uint) (SmsTransactionListResponse, error) {
	endpoint := fmt.Sprintf("/sms/%v/transactions?page=%v&limit=%v", serviceId, page, limit)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return SmsTransactionListResponse{}, err
	}
	var transactionList SmsTransactionListResponse
	return transactionList, json.Unmarshal(response, &transactionList)
}

func (s Sms) GetTransactionDetails(serviceId, transactionId uint) (SmsTransactionDetailsResponse, error) {
	endpoint := fmt.Sprintf("/sms/%v/transactions/%v", serviceId, transactionId)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return SmsTransactionDetailsResponse{}, err
	}
	var transactionDetails SmsTransactionDetailsResponse
	return transactionDetails, json.Unmarshal(response, &transactionDetails)
}

func (s Sms) GetServiceNumbers(serviceId, page, limit uint) (NumberListResponse, error) {
	endpoint := fmt.Sprintf("/sms/%v/numbers?page=%v&limit=%v", serviceId, page, limit)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return NumberListResponse{}, err
	}
	var numberListResponse NumberListResponse
	return numberListResponse, json.Unmarshal(response, &numberListResponse)
}

func (s Sms) GetServiceNumberDetails(serviceId uint, number int64) (NumberDetailsResponse, error) {
	endpoint := fmt.Sprintf("/sms/%v/numbers/%v", serviceId, number)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return NumberDetailsResponse{}, err
	}
	var numberDetails NumberDetailsResponse
	return numberDetails, json.Unmarshal(response, &numberDetails)
}

func (s Sms) GetNumbers(page, limit uint) (NumberListResponse, error) {
	endpoint := fmt.Sprintf("/sms/numbers?page=%v&limit=%v", page, limit)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return NumberListResponse{}, err
	}
	var numberListResponse NumberListResponse
	return numberListResponse, json.Unmarshal(response, &numberListResponse)
}

func (s Sms) GetNumberDetails(number int64) (NumberDetailsResponse, error) {
	endpoint := fmt.Sprintf("/sms/numbers/%v", number)
	response, err := s.restClient.sendGetRequest(endpoint)
	if err != nil {
		return NumberDetailsResponse{}, err
	}
	var numberDetails NumberDetailsResponse
	return numberDetails, json.Unmarshal(response, &numberDetails)
}

func (s Sms) VerifyCode(serviceId uint, code string, number int64) (CodeVerificationResponse, error) {
	endpoint := fmt.Sprintf("/sms/%v", serviceId)
	response, err := s.restClient.sendPostRequest(endpoint, CodeVerifyRequest{Code: code, Number: number})
	if err != nil {
		return CodeVerificationResponse{}, err
	}
	var verificationResponse CodeVerificationResponse
	return verificationResponse, json.Unmarshal(response, &verificationResponse)
}
