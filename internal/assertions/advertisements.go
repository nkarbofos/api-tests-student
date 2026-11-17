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

	require.Failf(t, "объявление не найдено", "ожидали объявление с названием %q среди %d найденных", title, len(items))
}

const priceDelta = 0.001

// AdvertisementMatchesCreateRequest проверяет, что параметры объявления совпадают с переданными при создании
func AdvertisementMatchesCreateRequest(t *testing.T, advertisement models.Advertisement, request models.AdvertisementCreateRequest) {
	t.Helper()

	require.Equal(t, request.Title, advertisement.Title, "title не совпадает")
	require.Equal(t, request.Description, advertisement.Description, "description не совпадает")

	if request.Price != nil {
		require.InDelta(t, *request.Price, advertisement.Price, priceDelta, "price не совпадает")
	}

	if request.Quantity != nil {
		require.Equal(t, *request.Quantity, advertisement.Quantity, "quantity не совпадает")
	}

	require.Len(t, advertisement.Photos, len(request.Photos), "количество фотографий не совпадает")
}
