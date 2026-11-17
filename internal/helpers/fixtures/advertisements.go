package fixtures

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/helpers/photos"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	foreignUserFirstName     = "Autotest"
	foreignUserLastName      = "Fixture"
	userPasswordRandomLength = 13
	fixturePrice             = 100.0
	fixtureQuantity          = 1
	fixtureDescription       = "Advertisement created by test fixture"
)

// RegisterTestUser регистрирует отдельного тестового пользователя и возвращает токен с данными пользователя
func RegisterTestUser(t *testing.T) authModels.LoginOkResponse {
	t.Helper()

	registerRequest := authModels.RegisterRequest{
		Email:     fmt.Sprintf("autotest-%s@example.com", utils.RandomString(12)),
		FirstName: foreignUserFirstName,
		LastName:  foreignUserLastName,
		Password:  fmt.Sprintf("Aa1%s", utils.RandomString(userPasswordRandomLength)),
	}

	return auth.Register(t, registerRequest)
}

// CreateForeignAdvertisement регистрирует отдельного пользователя и создаёт от его имени объявление, возвращая модель и токен владельца
func CreateForeignAdvertisement(t *testing.T, title string) (advertisementModels.Advertisement, string) {
	t.Helper()

	loginData := RegisterTestUser(t)

	price := fixturePrice
	quantity := fixtureQuantity
	createRequest := advertisementModels.AdvertisementCreateRequest{
		Title:       title,
		Description: fixtureDescription,
		Price:       &price,
		Quantity:    &quantity,
		Photos:      photos.CreateTempPNGs(t, 1),
	}

	advertisement := advertisementModels.Advertisement{}
	t.Cleanup(func() {
		if advertisement.Id != "" {
			advertisements.DeleteAdvertisementIfExists(t, loginData.Token, advertisement.Id)
		}
	})

	advertisement, _ = advertisements.CreateAdvertisement(t, loginData.Token, createRequest, http.StatusCreated)

	return advertisement, loginData.Token
}
