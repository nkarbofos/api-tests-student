package advertisementCreateNegative

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/apierrors"
	"api-tests-template/internal/helpers/photos"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	titleRandomLength       = 8
	validDescription        = "Объявление для проверки валидации"
	validPrice              = 100.0
	validQuantity           = 1
	titleMaxLength          = 50
	excessTitleLength       = 1
	validPhotoCount         = 1
	excessPhotoCount        = 4
	invalidToken            = "invalid_token"
	emptyTitleMessage       = "title is required"
	emptyDescriptionMessage = "description is required"
	priceRangeMessage       = "price must be between 0 and 1000000 inclusive"
	quantityRangeMessage    = "quantity must be between 1 and 100 inclusive"
	noPhotosMessage         = "at least 1 photo is required"
	tooManyPhotosMessage    = "maximum 3 photos allowed"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

type validationCase struct {
	name            string
	request         advertisementModels.AdvertisementCreateRequest
	expectedMessage string
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

// validAdvertisementRequest возвращает корректный запрос на создание объявления, от которого строятся негативные кейсы
func validAdvertisementRequest(t *testing.T) advertisementModels.AdvertisementCreateRequest {
	t.Helper()

	price := validPrice
	quantity := validQuantity

	return advertisementModels.AdvertisementCreateRequest{
		Title:       fmt.Sprintf("Autotest negative %s", utils.RandomString(titleRandomLength)),
		Description: validDescription,
		Price:       &price,
		Quantity:    &quantity,
		Photos:      photos.CreateTempPNGs(t, validPhotoCount),
	}
}

// requestWithMutation возвращает корректный запрос на создание объявления с применённой мутацией
func requestWithMutation(t *testing.T, mutate func(*advertisementModels.AdvertisementCreateRequest)) advertisementModels.AdvertisementCreateRequest {
	t.Helper()

	request := validAdvertisementRequest(t)
	mutate(&request)

	return request
}

func (s *TestSuite) TestCreateAdvertisementValidation() {
	base.Precondition("Подготовка набора запросов с некорректными параметрами объявления")

	negativePrice := -1.0
	tooBigPrice := 1000001.0
	negativeQuantity := -1

	cases := []validationCase{
		{
			name: "без обязательных параметров",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Title = ""
				request.Description = ""
				request.Photos = nil
			}),
			expectedMessage: emptyTitleMessage,
		},
		{
			name: "пустой title",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Title = ""
			}),
			expectedMessage: emptyTitleMessage,
		},
		{
			name: "пустой description",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Description = ""
			}),
			expectedMessage: emptyDescriptionMessage,
		},
		{
			name: "слишком длинный title",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Title = strings.Repeat("a", titleMaxLength+excessTitleLength)
			}),
			expectedMessage: fmt.Sprintf("title must not exceed %d characters provided %d", titleMaxLength, titleMaxLength+excessTitleLength),
		},
		{
			name: "отрицательная цена",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Price = &negativePrice
			}),
			expectedMessage: priceRangeMessage,
		},
		{
			name: "цена выше максимальной",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Price = &tooBigPrice
			}),
			expectedMessage: priceRangeMessage,
		},
		{
			name: "отрицательное количество",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Quantity = &negativeQuantity
			}),
			expectedMessage: quantityRangeMessage,
		},
		{
			name: "без фотографий",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Photos = nil
			}),
			expectedMessage: noPhotosMessage,
		},
		{
			name: "больше трёх фотографий",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Photos = photos.CreateTempPNGs(s.T(), excessPhotoCount)
			}),
			expectedMessage: tooManyPhotosMessage,
		},
	}

	for _, testCase := range cases {
		s.Run(testCase.name, func() {
			base.Step(fmt.Sprintf("POST /advertisement: %s", testCase.name))
			_, body := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, testCase.request, http.StatusBadRequest)
			assertions.ErrorResponseMatches(s.T(), body, apierrors.InvalidRequestCode, testCase.expectedMessage)
			base.Check("Получены ожидаемые код и текст ошибки")
		})
	}
}

func (s *TestSuite) TestCreateAdvertisementWithoutToken() {
	base.Precondition("Подготовка корректного запроса на создание объявления")
	base.Step("POST /advertisement без заголовка Authorization")
	_, body := advertisements.CreateAdvertisement(s.T(), "", validAdvertisementRequest(s.T()), http.StatusUnauthorized)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
	base.Check("Получен 401 Unauthorized с сообщением об отсутствии авторизации")
}

func (s *TestSuite) TestCreateAdvertisementWithInvalidToken() {
	base.Precondition("Подготовка корректного запроса на создание объявления")
	base.Step("POST /advertisement с неверным токеном")
	_, body := advertisements.CreateAdvertisement(s.T(), invalidToken, validAdvertisementRequest(s.T()), http.StatusUnauthorized)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.InvalidAuthorizationMessage)
	base.Check("Получен 401 Unauthorized с сообщением о неверном токене")
}
