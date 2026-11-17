package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

const favoritesPageSize = "50"

// HttpGetFavorites выполняет GET-запрос для получения избранного авторизованного пользователя
func HttpGetFavorites(t *testing.T, token, search string) *http.Response {
	apiTest := apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", "1").
		Query("page_size", favoritesPageSize)

	if search != "" {
		apiTest = apiTest.Query("search", search)
	}

	return apiTest.Expect(t).End().Response
}

// HttpAddFavorite выполняет POST-запрос для добавления объявления в избранное
func HttpAddFavorite(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(models.AddFavoriteRequest{AdvertisementId: advertisementId}).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite выполняет DELETE-запрос для удаления объявления из избранного
func HttpDeleteFavorite(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf(path.FavoritePath, advertisementId)).
		Expect(t).
		End().Response
}
