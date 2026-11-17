package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное, проверяет ожидаемый статус-код и возвращает модель и тело ответа
func AddFavorite(t *testing.T, token, advertisementId string, expectedStatusCode int) (models.FavoriteItem, string) {
	raw := favorites.HttpAddFavorite(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	favorite := models.FavoriteItem{}
	if expectedStatusCode == http.StatusCreated {
		require.NoErrorf(t, json.Unmarshal(body, &favorite), "не удалось разобрать тело ответа: %s", string(body))
	}

	return favorite, string(body)
}

// GetFavorites возвращает избранное пользователя с фильтром по названию объявления, проверяет ожидаемый статус-код
func GetFavorites(t *testing.T, token, search string, expectedStatusCode int) (models.FavoriteListResponse, string) {
	raw := favorites.HttpGetFavorites(t, token, search)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	favoritesList := models.FavoriteListResponse{}
	if expectedStatusCode == http.StatusOK {
		require.NoErrorf(t, json.Unmarshal(body, &favoritesList), "не удалось разобрать тело ответа: %s", string(body))
	}

	return favoritesList, string(body)
}

// DeleteFavorite удаляет объявление из избранного, проверяет ожидаемый статус-код и возвращает тело ответа
func DeleteFavorite(t *testing.T, token, advertisementId string, expectedStatusCode int) string {
	raw := favorites.HttpDeleteFavorite(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	return string(body)
}
