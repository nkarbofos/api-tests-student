package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// ErrorResponseMatches проверяет код и текст ошибки в теле ответа API
func ErrorResponseMatches(t *testing.T, body, expectedError, expectedMessage string) {
	t.Helper()

	require.Equal(t, expectedError, gjson.Get(body, "error").String(), "некорректный код ошибки, тело ответа: %s", body)
	require.Equal(t, expectedMessage, gjson.Get(body, "message").String(), "некорректный текст ошибки, тело ответа: %s", body)
}
