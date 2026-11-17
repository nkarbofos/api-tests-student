package auth

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "username is empty")
	require.NotEmpty(t, password, "password is empty")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// Register регистрирует нового пользователя и возвращает токен с данными пользователя, ожидает 201
func Register(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	registerRawResponse := auth.HttpPostAuthRegister(t, request)
	defer registerRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(registerRawResponse.Body)
	require.NoError(t, errResponse)

	registerResponse := models.LoginOkResponse{}
	require.NoErrorf(t, json.Unmarshal(respArray, &registerResponse), "failed to parse response body: %s", string(respArray))
	require.NotEmpty(t, registerResponse.Token, "registered user token is empty")

	return registerResponse
}
