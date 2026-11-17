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

	base.Precondition("Authorizing user with credentials from environment variables")
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
	description := fmt.Sprintf("advertisement %s did not appear in favorites", advertisementId)
	if !expectedPresent {
		description = fmt.Sprintf("advertisement %s did not disappear from favorites", advertisementId)
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

	base.Precondition("Creating advertisement as a separate test user")
	foreignAdvertisement, ownerToken := fixtures.CreateForeignAdvertisement(s.T(), title)
	base.Check(fmt.Sprintf("Fixture advertisement created: id=%s", foreignAdvertisement.Id))

	s.Run("Add advertisement to favorites", func() {
		base.Step(fmt.Sprintf("POST /favorites with advertisement_id=%s", foreignAdvertisement.Id))
		favorite, body := favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusCreated)

		require.Equal(s.T(), foreignAdvertisement.Id, favorite.Advertisement.Id, "advertisement_id in response does not match; response body: %s", body)
		require.NotEmpty(s.T(), favorite.FavoriteId, "favorite_id is empty; response body: %s", body)
		base.Check("Advertisement added to favorites")
	})

	s.Run("Advertisement appears in favorites", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", title))
		items := s.waitForFavoriteState(title, foreignAdvertisement.Id, true)
		assertions.FavoritesContainAdvertisement(s.T(), items, foreignAdvertisement.Id)
		base.Check("Advertisement found in favorites")
	})

	s.Run("Remove advertisement from favorites", func() {
		base.Step(fmt.Sprintf("DELETE /favorites/%s", foreignAdvertisement.Id))
		favoritesManager.DeleteFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusNoContent)
		base.Check("Delete favorite endpoint returned 204 No Content")
	})

	s.Run("Removed advertisement disappears from favorites", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", title))
		items := s.waitForFavoriteState(title, foreignAdvertisement.Id, false)
		assertions.FavoritesDoNotContainAdvertisement(s.T(), items, foreignAdvertisement.Id)
		base.Check("Advertisement is absent from favorites")
	})

	s.Run("Add advertisement to favorites again", func() {
		base.Step(fmt.Sprintf("POST /favorites with advertisement_id=%s", foreignAdvertisement.Id))
		favoritesManager.AddFavorite(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusCreated)
		base.Check("Advertisement added to favorites again")
	})

	s.Run("Delete advertisement as owner", func() {
		base.Step(fmt.Sprintf("DELETE /advertisement?id=%s", foreignAdvertisement.Id))
		advertisements.DeleteAdvertisement(s.T(), ownerToken, foreignAdvertisement.Id, http.StatusNoContent)
		base.Check("Advertisement deleted by owner")
	})

	s.Run("Advertisement disappears from favorites after deletion", func() {
		base.Step(fmt.Sprintf("GET /advertisement?id=%s", foreignAdvertisement.Id))
		_, advertisementBody := advertisements.GetAdvertisement(s.T(), s.loginData.Token, foreignAdvertisement.Id, http.StatusNotFound)
		assertions.ErrorResponseMatches(s.T(), advertisementBody, apierrors.NotFoundCode, apierrors.AdvertisementNotFoundMessage)
		base.Check("Deleted advertisement is unavailable via GET /advertisement")

		base.Step(fmt.Sprintf("GET /favorites?search=%s", title))
		items := s.waitForFavoriteState(title, foreignAdvertisement.Id, false)
		assertions.FavoritesDoNotContainAdvertisement(s.T(), items, foreignAdvertisement.Id)
		base.Check("Advertisement is absent from favorites after deletion")
	})
}
