package favorites

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/apierrors"
	"api-tests-template/internal/helpers/fixtures"
	"api-tests-template/internal/helpers/wait"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favoritesManager "api-tests-template/internal/managers/favorites"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	favoriteTitleRandomLength = 8
	favoritesWaitTimeout      = 10 * time.Second
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// containsFavorite проверяет наличие объявления в списке избранного без прерывания теста
func containsFavorite(items []favoriteModels.FavoriteItem, advertisementId string) bool {
	for _, item := range items {
		if item.Advertisement.Id == advertisementId {
			return true
		}
	}
	return false
}

// waitForFavoriteState ждёт, пока объявление появится или исчезнет в избранном, и возвращает последний полученный список
func (s *TestSuite) waitForFavoriteState(title, advertisementId string, expectedPresent bool) []favoriteModels.FavoriteItem {
	var items []favoriteModels.FavoriteItem
	description := fmt.Sprintf("объявление %s не появилось в избранном", advertisementId)
	if !expectedPresent {
		description = fmt.Sprintf("объявление %s не исчезло из избранного", advertisementId)
	}

	wait.Until(s.T(), favoritesWaitTimeout, description, func() bool {
		list, _ := favoritesManager.GetFavorites(s.T(), s.loginData.Token, title, http.StatusOK)
		items = list.Items
		return containsFavorite(list.Items, advertisementId) == expectedPresent
	})

	return items
}

func (s *TestSuite) TestFavoritesFullCycle() {
	title := fmt.Sprintf("Autotest favorite %s", utils.RandomString(favoriteTitleRandomLength))

	base.Precondition("Создание объявления от имени отдельного тестового пользователя")
	foreignAdvertisement, ownerToken := fixtures.CreateForeignAdvertisement(s.T(), title)
	base.Check(fmt.Sprintf("Объявление фикстуры создано: id=%s", foreignAdvertisement.Id))

	s.Run("Объявление добавляется в избранное", func() {
		base.Step(fmt.Sprintf("POST /favorites с advertisement_id=%s", foreignAdvertisement.Id))
		favorite, body := favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusCreated)

		require.Equal(s.T(), foreignAdvertisement.Id, favorite.Advertisement.Id, "advertisement_id в ответе не совпадает; тело ответа: %s", body)
		require.NotEmpty(s.T(), favorite.FavoriteId, "favorite_id пуст; тело ответа: %s", body)
		base.Check("Объявление добавлено в избранное")
	})

	s.Run("Добавленное объявление появляется в списке избранного", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", title))
		items := s.waitForFavoriteState(title, foreignAdvertisement.Id, true)
		assertions.FavoritesContainAdvertisement(s.T(), items, foreignAdvertisement.Id)
		base.Check("Объявление найдено в списке избранного")
	})

	s.Run("Объявление удаляется из избранного", func() {
		base.Step(fmt.Sprintf("DELETE /favorites/%s", foreignAdvertisement.Id))
		favoritesManager.DeleteFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusNoContent)
		base.Check("Ручка удаления из избранного вернула 204 No Content")
	})

	s.Run("Удалённое объявление пропадает из списка избранного", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", title))
		items := s.waitForFavoriteState(title, foreignAdvertisement.Id, false)
		assertions.FavoritesDoNotContainAdvertisement(s.T(), items, foreignAdvertisement.Id)
		base.Check("Объявление отсутствует в списке избранного")
	})

	s.Run("Повторное добавление объявления в избранное", func() {
		base.Step(fmt.Sprintf("POST /favorites с advertisement_id=%s", foreignAdvertisement.Id))
		favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusCreated)
		base.Check("Объявление повторно добавлено в избранное")
	})

	s.Run("Удаление объявления владельцем", func() {
		base.Step(fmt.Sprintf("DELETE /advertisement?id=%s", foreignAdvertisement.Id))
		advertisements.DeleteAdvertisement(s.T(), ownerToken, foreignAdvertisement.Id, http.StatusNoContent)
		base.Check("Объявление удалено владельцем")
	})

	s.Run("Объявление исчезает из избранного после удаления самого объявления", func() {
		base.Step(fmt.Sprintf("GET /advertisement?id=%s", foreignAdvertisement.Id))
		_, advertisementBody := advertisements.GetAdvertisement(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusNotFound)
		assertions.ErrorResponseMatches(s.T(), advertisementBody, apierrors.NotFoundCode, apierrors.AdvertisementNotFoundMessage)
		base.Check("Удалённое объявление недоступно через GET /advertisement")

		base.Step(fmt.Sprintf("GET /favorites?search=%s", title))
		items := s.waitForFavoriteState(title, foreignAdvertisement.Id, false)
		assertions.FavoritesDoNotContainAdvertisement(s.T(), items, foreignAdvertisement.Id)
		base.Check("Объявление отсутствует в избранном после удаления объявления")
	})
}
