package api_runner

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const (
	requestAttempts = 5
	baseRetryDelay  = 300 * time.Millisecond
)

// retryTransport повторяет запрос при временных ответах сервера 429/502/503/504, чтобы снизить нестабильность тестов
type retryTransport struct {
	transport http.RoundTripper
}

// RoundTrip выполняет запрос и повторяет его при временных ошибках сервера
func (retry *retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	body, err := readRequestBody(request)
	if err != nil {
		return nil, err
	}

	var response *http.Response

	for attempt := 1; attempt <= requestAttempts; attempt++ {
		attemptRequest := request.Clone(request.Context())
		if body != nil {
			attemptRequest.Body = io.NopCloser(bytes.NewReader(body))
		}

		response, err = retry.transport.RoundTrip(attemptRequest)
		if err != nil {
			return response, err
		}

		if !isRetryableStatusCode(response.StatusCode) || attempt == requestAttempts {
			return response, nil
		}

		response.Body.Close()
		time.Sleep(retryDelayForAttempt(attempt))
	}

	return response, err
}

// retryDelayForAttempt возвращает экспоненциальную задержку перед следующей попыткой
func retryDelayForAttempt(attempt int) time.Duration {
	return baseRetryDelay * time.Duration(1<<(attempt-1))
}

// readRequestBody читает тело запроса, чтобы его можно было повторно отправить при retry
func readRequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}
	defer request.Body.Close()

	return io.ReadAll(request.Body)
}

// isRetryableStatusCode сообщает, является ли статус-код временной ошибкой сервера
func isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
