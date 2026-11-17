package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// FavoritesContainAdvertisement проверяет, что объявление есть в списке избранного
func FavoritesContainAdvertisement(t *testing.T, items []models.FavoriteItem, advertisementId string) {
	t.Helper()

	for _, item := range items {
		if item.Advertisement.Id == advertisementId {
			return
		}
	}

	require.Failf(t, "advertisement not found in favorites", "expected advertisement with id %q among %d favorite entries", advertisementId, len(items))
}

// FavoritesDoNotContainAdvertisement проверяет, что объявления нет в списке избранного
func FavoritesDoNotContainAdvertisement(t *testing.T, items []models.FavoriteItem, advertisementId string) {
	t.Helper()

	for _, item := range items {
		if item.Advertisement.Id == advertisementId {
			require.Failf(t, "advertisement is still in favorites", "did not expect advertisement with id %q among %d favorite entries", advertisementId, len(items))
		}
	}
}
