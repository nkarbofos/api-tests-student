package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "advertisement not found", "expected advertisement with title %q among %d found", title, len(items))
}

const priceDelta = 0.001

// AdvertisementMatchesCreateRequest проверяет, что параметры объявления совпадают с переданными при создании
func AdvertisementMatchesCreateRequest(t *testing.T, advertisement models.Advertisement, request models.AdvertisementCreateRequest) {
	t.Helper()

	require.Equal(t, request.Title, advertisement.Title, "title does not match")
	require.Equal(t, request.Description, advertisement.Description, "description does not match")

	if request.Price != nil {
		require.InDelta(t, *request.Price, advertisement.Price, priceDelta, "price does not match")
	}

	if request.Quantity != nil {
		require.Equal(t, *request.Quantity, advertisement.Quantity, "quantity does not match")
	}

	require.Len(t, advertisement.Photos, len(request.Photos), "photos count does not match")
}
