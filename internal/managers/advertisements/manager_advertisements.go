package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// GetAdvertisementsList возвращает список объявлений, найденных по вхождению подстроки в название
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// CreateAdvertisement создаёт объявление, проверяет ожидаемый статус-код и возвращает модель и тело ответа
func CreateAdvertisement(t *testing.T, token string, request models.AdvertisementCreateRequest, expectedStatusCode int) (models.Advertisement, string) {
	raw := advertisements.HttpCreateAdvertisement(t, token, request)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	advertisement := models.Advertisement{}
	if expectedStatusCode == http.StatusCreated {
		require.NoErrorf(t, json.Unmarshal(body, &advertisement), "не удалось разобрать тело ответа: %s", string(body))
	}

	return advertisement, string(body)
}

// GetAdvertisement получает объявление по id, проверяет ожидаемый статус-код и возвращает модель и тело ответа
func GetAdvertisement(t *testing.T, token, advertisementId string, expectedStatusCode int) (models.Advertisement, string) {
	raw := advertisements.HttpGetAdvertisement(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	advertisement := models.Advertisement{}
	if expectedStatusCode == http.StatusOK {
		require.NoErrorf(t, json.Unmarshal(body, &advertisement), "не удалось разобрать тело ответа: %s", string(body))
	}

	return advertisement, string(body)
}

// GetAdvertisementPhotos получает фотографии объявления, проверяет ожидаемый статус-код и возвращает модели и тело ответа
func GetAdvertisementPhotos(t *testing.T, token, advertisementId string, expectedStatusCode int) ([]models.Photo, string) {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	photos := []models.Photo{}
	if expectedStatusCode == http.StatusOK {
		require.NoErrorf(t, json.Unmarshal(body, &photos), "не удалось разобрать тело ответа: %s", string(body))
	}

	return photos, string(body)
}

// DeleteAdvertisement удаляет объявление, проверяет ожидаемый статус-код и возвращает тело ответа
func DeleteAdvertisement(t *testing.T, token, advertisementId string, expectedStatusCode int) string {
	raw := advertisements.HttpDeleteAdvertisement(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	return string(body)
}

// DeleteAdvertisementIfExists удаляет объявление при очистке тестовых данных и допускает, что оно уже удалено
func DeleteAdvertisementIfExists(t *testing.T, token, advertisementId string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Containsf(t, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, raw.StatusCode,
		"ожидали 204, 403 или 404 при удалении объявления, получен %d, тело ответа: %s", raw.StatusCode, string(body))
}
