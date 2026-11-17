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
	validDescription        = "Advertisement for validation checks"
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

	base.Precondition("Authorizing user with credentials from environment variables")
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
	base.Precondition("Preparing requests with invalid advertisement parameters")

	negativePrice := -1.0
	tooBigPrice := 1000001.0
	negativeQuantity := -1

	cases := []validationCase{
		{
			name: "missing required parameters",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Title = ""
				request.Description = ""
				request.Photos = nil
			}),
			expectedMessage: emptyTitleMessage,
		},
		{
			name: "empty title",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Title = ""
			}),
			expectedMessage: emptyTitleMessage,
		},
		{
			name: "empty description",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Description = ""
			}),
			expectedMessage: emptyDescriptionMessage,
		},
		{
			name: "title too long",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Title = strings.Repeat("a", titleMaxLength+excessTitleLength)
			}),
			expectedMessage: fmt.Sprintf("title must not exceed %d characters provided %d", titleMaxLength, titleMaxLength+excessTitleLength),
		},
		{
			name: "negative price",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Price = &negativePrice
			}),
			expectedMessage: priceRangeMessage,
		},
		{
			name: "price above maximum",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Price = &tooBigPrice
			}),
			expectedMessage: priceRangeMessage,
		},
		{
			name: "negative quantity",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Quantity = &negativeQuantity
			}),
			expectedMessage: quantityRangeMessage,
		},
		{
			name: "no photos",
			request: requestWithMutation(s.T(), func(request *advertisementModels.AdvertisementCreateRequest) {
				request.Photos = nil
			}),
			expectedMessage: noPhotosMessage,
		},
		{
			name: "more than three photos",
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
			base.Check("Expected error code and message received")
		})
	}
}

func (s *TestSuite) TestCreateAdvertisementWithoutToken() {
	base.Precondition("Preparing valid advertisement create request")
	base.Step("POST /advertisement without Authorization header")
	_, body := advertisements.CreateAdvertisement(s.T(), "", validAdvertisementRequest(s.T()), http.StatusUnauthorized)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
	base.Check("Received 401 Unauthorized with missing authorization message")
}

func (s *TestSuite) TestCreateAdvertisementWithInvalidToken() {
	base.Precondition("Preparing valid advertisement create request")
	base.Step("POST /advertisement with invalid token")
	_, body := advertisements.CreateAdvertisement(s.T(), invalidToken, validAdvertisementRequest(s.T()), http.StatusUnauthorized)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.InvalidAuthorizationMessage)
	base.Check("Received 401 Unauthorized with invalid token message")
}
