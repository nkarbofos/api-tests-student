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
	advertisementDescription = "Advertisement created in positive scenario"
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

	base.Precondition("Registering a separate test user for the scenario")
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

	base.Precondition("Creating advertisement with all possible parameters")
	base.Step("POST /advertisement with title, description, price, quantity and photos")
	created, createBody := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request, http.StatusCreated)

	s.T().Cleanup(func() {
		if created.Id == "" {
			return
		}
		base.Precondition("Cleanup: deleting created advertisement")
		advertisements.DeleteAdvertisementIfExists(s.T(), s.loginData.Token, created.Id)
	})

	require.NotEmpty(s.T(), created.Id, "created advertisement id is empty; response body: %s", createBody)
	require.Equal(s.T(), s.loginData.User.Id, created.UserId, "created advertisement user_id does not match current user")
	require.Len(s.T(), created.Photos, photoCount, "photos count in create response does not match")
	base.Check(fmt.Sprintf("Advertisement created: id=%s", created.Id))

	s.Run("Advertisement is available by id", func() {
		base.Step(fmt.Sprintf("GET /advertisement?id=%s", created.Id))
		fetched, fetchedBody := advertisements.GetAdvertisement(s.T(), s.loginData.Token, created.Id, http.StatusOK)

		require.Equal(s.T(), created.Id, fetched.Id, "fetched advertisement id does not match; response body: %s", fetchedBody)
		require.Equal(s.T(), s.loginData.User.Id, fetched.UserId, "fetched advertisement user_id does not match")
		assertions.AdvertisementMatchesCreateRequest(s.T(), fetched, request)

		for index, photo := range fetched.Photos {
			require.NotEmpty(s.T(), photo.Url, "photo %d url is empty", index)
		}

		base.Check("Advertisement parameters match the create request")
	})

	s.Run("Advertisement is found by full title", func() {
		base.Step(fmt.Sprintf("GET /advertisements?search=%s", request.Title))

		wait.Until(s.T(), searchWaitTimeout, fmt.Sprintf("advertisement %s did not appear in search by title %q", created.Id, request.Title), func() bool {
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, request.Title)
			for _, item := range list.Items {
				if item.Id == created.Id && item.Title == request.Title {
					return true
				}
			}
			return false
		})

		base.Check("Advertisement found in search by full title")
	})

	s.Run("Photos are available", func() {
		base.Step(fmt.Sprintf("GET /advertisements/%s/photos", created.Id))
		advertisementPhotos, photosBody := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.Id, http.StatusOK)
		require.Len(s.T(), advertisementPhotos, photoCount, "photos count does not match; response body: %s", photosBody)

		for index, photo := range advertisementPhotos {
			require.NotEmpty(s.T(), photo.Id, "photo %d id is empty", index)
			require.NotEmpty(s.T(), photo.Url, "photo %d url is empty", index)
		}

		base.Check("All uploaded photos are available on the server")
	})
}
