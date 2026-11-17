package advertisements

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет GET-запрос для поиска объявлений по названию
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpCreateAdvertisement выполняет POST-запрос для создания объявления с фотографиями
func HttpCreateAdvertisement(t *testing.T, token string, request models.AdvertisementCreateRequest) *http.Response {
	apiTest := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description)

	if request.Price != nil {
		apiTest = apiTest.MultipartFormData("price", strconv.FormatFloat(*request.Price, 'f', -1, 64))
	}

	if request.Quantity != nil {
		apiTest = apiTest.MultipartFormData("quantity", strconv.Itoa(*request.Quantity))
	}

	if len(request.Photos) > 0 {
		apiTest = apiTest.MultipartFile("photos", request.Photos...)
	}

	return apiTest.Expect(t).End().Response
}

// HttpGetAdvertisement выполняет GET-запрос для получения объявления по id
func HttpGetAdvertisement(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", advertisementId).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос для удаления объявления по id
func HttpDeleteAdvertisement(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", advertisementId).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос для получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(fmt.Sprintf(path.AdvertisementPhotosPath, advertisementId)).
		Expect(t).
		End().Response
}
