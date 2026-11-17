package favoritesNegative

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/apierrors"
	"api-tests-template/internal/helpers/fixtures"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favoritesManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	nonExistentAdvertisementId = "00000000-0000-0000-0000-000000000001"
	ownAdvertisementTitle      = "Autotest own favorite "
	titleRandomLength          = 8
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

func (s *TestSuite) TestAddFavoriteWithNonExistentAdvertisement() {
	base.Precondition("Подготовка идентификатора несуществующего объявления")
	base.Step(fmt.Sprintf("POST /favorites с несуществующим advertisement_id=%s", nonExistentAdvertisementId))
	_, body := favoritesManager.AddFavorite(s.T(), s.loginData.Token, nonExistentAdvertisementId, http.StatusNotFound)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.NotFoundCode, apierrors.AdvertisementNotFoundMessage)
	base.Check("Получен 404 Not Found с сообщением об отсутствии объявления")
}

func (s *TestSuite) TestDeleteFavoriteWithNonExistentAdvertisement() {
	base.Precondition("Подготовка идентификатора несуществующего объявления")
	base.Step(fmt.Sprintf("DELETE /favorites/%s", nonExistentAdvertisementId))
	body := favoritesManager.DeleteFavorite(s.T(), s.loginData.Token, nonExistentAdvertisementId, http.StatusNotFound)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.NotFoundCode, apierrors.FavoriteNotFoundMessage)
	base.Check("Получен 404 Not Found с сообщением об отсутствии записи в избранном")
}

func (s *TestSuite) TestAddFavoriteDuplicate() {
	title := fmt.Sprintf("Autotest favorite %s", utils.RandomString(titleRandomLength))

	base.Precondition("Создание объявления от имени отдельного тестового пользователя")
	foreignAdvertisement, _ := fixtures.CreateForeignAdvertisement(s.T(), title)

	base.Step(fmt.Sprintf("Первое добавление объявления %s в избранное", foreignAdvertisement.Id))
	favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusCreated)
	base.Check("Объявление добавлено в избранное")

	base.Step(fmt.Sprintf("Повторное добавление объявления %s в избранное", foreignAdvertisement.Id))
	_, body := favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusConflict)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.ConflictCode, apierrors.AdvertisementAlreadyFavoriteMessage)
	base.Check("Получен 409 Conflict с сообщением о том, что объявление уже в избранном")
}

func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	title := fmt.Sprintf("%s%s", ownAdvertisementTitle, utils.RandomString(titleRandomLength))

	base.Precondition("Создание объявления от имени отдельного тестового пользователя")
	ownAdvertisement, ownerToken := fixtures.CreateForeignAdvertisement(s.T(), title)
	base.Check(fmt.Sprintf("Собственное объявление пользователя создано: id=%s", ownAdvertisement.Id))

	base.Step(fmt.Sprintf("POST /favorites со своим объявлением %s", ownAdvertisement.Id))
	_, errorBody := favoritesManager.AddFavorite(s.T(), ownerToken, ownAdvertisement.Id, http.StatusForbidden)
	assertions.ErrorResponseMatches(s.T(), errorBody, apierrors.ForbiddenCode, apierrors.OwnAdvertisementFavoriteMessage)
	base.Check("Получен 403 Forbidden с сообщением о запрете добавления своего объявления")
}

func (s *TestSuite) TestFavoritesWithoutAuthorization() {
	base.Precondition("Подготовка запросов к ручкам избранного без авторизации")

	s.Run("GET /favorites без авторизации", func() {
		base.Step("GET /favorites без заголовка Authorization")
		_, body := favoritesManager.GetFavorites(s.T(), "", "", http.StatusUnauthorized)
		assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
		base.Check("Получен 401 Unauthorized")
	})

	s.Run("POST /favorites без авторизации", func() {
		base.Step("POST /favorites без заголовка Authorization")
		_, body := favoritesManager.AddFavorite(s.T(), "", nonExistentAdvertisementId, http.StatusUnauthorized)
		assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
		base.Check("Получен 401 Unauthorized")
	})

	s.Run("DELETE /favorites без авторизации", func() {
		base.Step("DELETE /favorites без заголовка Authorization")
		body := favoritesManager.DeleteFavorite(s.T(), "", nonExistentAdvertisementId, http.StatusUnauthorized)
		assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
		base.Check("Получен 401 Unauthorized")
	})
}
