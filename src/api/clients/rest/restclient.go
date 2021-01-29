package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

// Mock struct
type Mock struct {
	URL        string
	HTTPMethod string
	Response   *http.Response
	Err        error
}

func getMockID(httpMethod string, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

// StartMockup func
func StartMockup() {
	enabledMocks = true
}

// FlushMockups func
func FlushMockups() {
	mocks = make(map[string]*Mock)
}

// StopMockup func
func StopMockup() {
	enabledMocks = false
}

// AddMockup func
func AddMockup(mock Mock) {
	mocks[getMockID(mock.HTTPMethod, mock.URL)] = &mock
}

// Post func
func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[getMockID(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("no mockup found for given request")
		}
		return mock.Response, mock.Err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}
