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

	base.Precondition("Authorizing user with credentials from environment variables")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestAddFavoriteWithNonExistentAdvertisement() {
	base.Precondition("Preparing non-existent advertisement id")
	base.Step(fmt.Sprintf("POST /favorites with non-existent advertisement_id=%s", nonExistentAdvertisementId))
	_, body := favoritesManager.AddFavorite(s.T(), s.loginData.Token, nonExistentAdvertisementId, http.StatusNotFound)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.NotFoundCode, apierrors.AdvertisementNotFoundMessage)
	base.Check("Received 404 Not Found with advertisement not found message")
}

func (s *TestSuite) TestDeleteFavoriteWithNonExistentAdvertisement() {
	base.Precondition("Preparing non-existent advertisement id")
	base.Step(fmt.Sprintf("DELETE /favorites/%s", nonExistentAdvertisementId))
	body := favoritesManager.DeleteFavorite(s.T(), s.loginData.Token, nonExistentAdvertisementId, http.StatusNotFound)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.NotFoundCode, apierrors.FavoriteNotFoundMessage)
	base.Check("Received 404 Not Found with favorite entry not found message")
}

func (s *TestSuite) TestAddFavoriteDuplicate() {
	title := fmt.Sprintf("Autotest favorite %s", utils.RandomString(titleRandomLength))

	base.Precondition("Creating advertisement as a separate test user")
	foreignAdvertisement, _ := fixtures.CreateForeignAdvertisement(s.T(), title)

	base.Step(fmt.Sprintf("First add of advertisement %s to favorites", foreignAdvertisement.Id))
	favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusCreated)
	base.Check("Advertisement added to favorites")

	base.Step(fmt.Sprintf("Repeated add of advertisement %s to favorites", foreignAdvertisement.Id))
	_, body := favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusConflict)
	assertions.ErrorResponseMatches(s.T(), body, apierrors.ConflictCode, apierrors.AdvertisementAlreadyFavoriteMessage)
	base.Check("Received 409 Conflict with already in favorites message")
}

func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	title := fmt.Sprintf("%s%s", ownAdvertisementTitle, utils.RandomString(titleRandomLength))

	base.Precondition("Creating advertisement as a separate test user")
	ownAdvertisement, ownerToken := fixtures.CreateForeignAdvertisement(s.T(), title)
	base.Check(fmt.Sprintf("User's own advertisement created: id=%s", ownAdvertisement.Id))

	base.Step(fmt.Sprintf("POST /favorites with own advertisement %s", ownAdvertisement.Id))
	_, errorBody := favoritesManager.AddFavorite(s.T(), ownerToken, ownAdvertisement.Id, http.StatusForbidden)
	assertions.ErrorResponseMatches(s.T(), errorBody, apierrors.ForbiddenCode, apierrors.OwnAdvertisementFavoriteMessage)
	base.Check("Received 403 Forbidden with own advertisement restriction message")
}

func (s *TestSuite) TestFavoritesWithoutAuthorization() {
	base.Precondition("Preparing unauthorized favorites requests")

	s.Run("GET /favorites unauthorized", func() {
		base.Step("GET /favorites without Authorization header")
		_, body := favoritesManager.GetFavorites(s.T(), "", "", http.StatusUnauthorized)
		assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
		base.Check("Received 401 Unauthorized")
	})

	s.Run("POST /favorites unauthorized", func() {
		base.Step("POST /favorites without Authorization header")
		_, body := favoritesManager.AddFavorite(s.T(), "", nonExistentAdvertisementId, http.StatusUnauthorized)
		assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
		base.Check("Received 401 Unauthorized")
	})

	s.Run("DELETE /favorites unauthorized", func() {
		base.Step("DELETE /favorites without Authorization header")
		body := favoritesManager.DeleteFavorite(s.T(), "", nonExistentAdvertisementId, http.StatusUnauthorized)
		assertions.ErrorResponseMatches(s.T(), body, apierrors.UnauthorizedCode, apierrors.AuthorizationHeaderRequiredMessage)
		base.Check("Received 401 Unauthorized")
	})
}
