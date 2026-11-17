package advertisementCreate

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/helpers/fixtures"
	"api-tests-template/internal/helpers/photos"
	"api-tests-template/internal/helpers/wait"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	advertisementDescription = "Объявление, созданное в позитивном сценарии"
	advertisementPrice       = 1234.5
	advertisementQuantity    = 7
	photoCount               = 2
	titleRandomLength        = 10
	searchWaitTimeout        = 10 * time.Second
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

	base.Precondition("Регистрация отдельного тестового пользователя для сценария")
	s.loginData = fixtures.RegisterTestUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestCreateAdvertisementWithAllParameters() {
	price := advertisementPrice
	quantity := advertisementQuantity
	request := advertisementModels.AdvertisementCreateRequest{
		Title:       fmt.Sprintf("Autotest ad %s", utils.RandomString(titleRandomLength)),
		Description: advertisementDescription,
		Price:       &price,
		Quantity:    &quantity,
		Photos:      photos.CreateTempPNGs(s.T(), photoCount),
	}

	base.Precondition("Создание объявления со всеми возможными параметрами")
	base.Step("POST /advertisement с title, description, price, quantity и фотографиями")
	created, createBody := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request, http.StatusCreated)

	s.T().Cleanup(func() {
		if created.Id == "" {
			return
		}
		base.Precondition("Очистка: удаление созданного объявления")
		advertisements.DeleteAdvertisementIfExists(s.T(), s.loginData.Token, created.Id)
	})

	require.NotEmpty(s.T(), created.Id, "id созданного объявления пуст; тело ответа: %s", createBody)
	require.Equal(s.T(), s.loginData.User.Id, created.UserId, "user_id созданного объявления не совпадает с текущим пользователем")
	require.Len(s.T(), created.Photos, photoCount, "количество фотографий в ответе на создание не совпадает")
	base.Check(fmt.Sprintf("Объявление создано: id=%s", created.Id))

	s.Run("Объявление доступно через GET /advertisement", func() {
		base.Step(fmt.Sprintf("GET /advertisement?id=%s", created.Id))
		fetched, fetchedBody := advertisements.GetAdvertisement(s.T(), s.loginData.Token, created.Id, http.StatusOK)

		require.Equal(s.T(), created.Id, fetched.Id, "id полученного объявления не совпадает; тело ответа: %s", fetchedBody)
		require.Equal(s.T(), s.loginData.User.Id, fetched.UserId, "user_id полученного объявления не совпадает")
		assertions.AdvertisementMatchesCreateRequest(s.T(), fetched, request)

		for index, photo := range fetched.Photos {
			require.NotEmpty(s.T(), photo.Url, "url фотографии %d пуст", index)
		}

		base.Check("Параметры объявления совпадают с переданными при создании")
	})

	s.Run("Объявление появляется в поиске по полному названию", func() {
		base.Step(fmt.Sprintf("GET /advertisements?search=%s", request.Title))

		wait.Until(s.T(), searchWaitTimeout, fmt.Sprintf("объявление %s не появилось в поиске по названию %q", created.Id, request.Title), func() bool {
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, request.Title)
			for _, item := range list.Items {
				if item.Id == created.Id && item.Title == request.Title {
					return true
				}
			}
			return false
		})

		base.Check("Объявление найдено в поиске по полному названию")
	})

	s.Run("Фотографии доступны через GET /advertisements/{id}/photos", func() {
		base.Step(fmt.Sprintf("GET /advertisements/%s/photos", created.Id))
		advertisementPhotos, photosBody := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.Id, http.StatusOK)
		require.Len(s.T(), advertisementPhotos, photoCount, "количество фотографий не совпадает; тело ответа: %s", photosBody)

		for index, photo := range advertisementPhotos {
			require.NotEmpty(s.T(), photo.Id, "id фотографии %d пуст", index)
			require.NotEmpty(s.T(), photo.Url, "url фотографии %d пуст", index)
		}

		base.Check("Все загруженные фотографии доступны на сервере")
	})
}
