package simpay

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	SimKeyHeader      = "X-SIM-KEY"
	SimPasswordHeader = "X-SIM-PASSWORD"
	SimVersionHeader  = "X-SIM-VERSION"
	SimPlatformHeader = "X-SIM-PLATFORM"
	BaseUrl           = "https://api.simpay.pl"
	ContentType       = "application/json"
)

type restClient struct {
	http.Client
}

func newRestClient(apiKey, simPassword string) restClient {
	return restClient{http.Client{Transport: interceptor{apikey: apiKey, simKey: simPassword, core: http.DefaultTransport}}}
}
func (r restClient) sendGetRequest(endpoint string) ([]byte, error) {
	resp, err := r.Get(BaseUrl + endpoint)
	return extractBody(err, resp)
}

func (r restClient) sendPostRequest(endpoint string, body interface{}) ([]byte, error) {
	marshaledBody, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	resp, err := r.Post(BaseUrl+endpoint, ContentType, bytes.NewBuffer(marshaledBody))
	return extractBody(err, resp)
}

func extractBody(err error, resp *http.Response) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	byteResponse, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return byteResponse, nil
}

type interceptor struct {
	apikey string
	simKey string
	core   http.RoundTripper
}

func (i interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set(SimKeyHeader, i.apikey)
	r.Header.Set(SimPasswordHeader, i.simKey)
	r.Header.Set(SimVersionHeader, "2.1.1")
	r.Header.Set(SimPlatformHeader, "GO")
	return i.core.RoundTrip(r)
}
