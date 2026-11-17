package models

type AdvertisementCreateRequest struct {
	Title       string
	Description string
	Price       *float64
	Quantity    *int
	Photos      []string
}
