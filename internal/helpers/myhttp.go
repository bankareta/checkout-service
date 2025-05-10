package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type ApiCall struct {
}

type ApiCallInterface interface {
	Send(ctx context.Context, url, method string, body []byte, headers map[string]string, timeout int64) ([]byte, http.Header, error)
}

// PerformAPICall function makes a dynamic API call based on the provided parameters
func (apiCall *ApiCall) Send(ctx context.Context, url, method string, body []byte, headers map[string]string, timeout int64) ([]byte, http.Header, error) {
	// Create a new HTTP client
	client := &http.Client{Timeout: time.Second * time.Duration(timeout)}
	logger := logrus.New()
	logger.Info(string(body))
	ctx.Value("refnum")
	// Create a new request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, nil, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
		}
		return responseBody, resp.Header, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
	}

	logger.Info(string(responseBody))

	return responseBody, resp.Header, nil
}

func GenerateHeader() http.Header {
	contentType := "application/json"
	headers := http.Header{
		"Content-Type": []string{contentType},
	}

	return headers
}

type ResponseLibV1[T any] struct {
	ResponseCode        string `json:"response_code"`
	ResponseRefnum      string `json:"response_refnum"`
	ResponseID          string `json:"response_id"`
	ResponseDescription string `json:"response_description"`
	ResponseData        T      `json:"response_data"`
}

type ResponseLibV2[T any] struct {
	ResponseCode   string      `json:"responseCode"`
	ResponseDesc   string      `json:"responseDesc"`
	ResponseData   T           `json:"responseData"`
	ResponseErrors interface{} `json:"responseErrors"`
}

type ResponseLibV3[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func IsMapStringString(data interface{}) (result map[string]string) {
	value, ok := data.(map[string]string)
	if ok {
		result = value
	}
	return
}
