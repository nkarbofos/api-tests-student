package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Advertisement struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserId      string  `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Photos      []Photo `json:"photos"`
}

type Photo struct {
	Id              string `json:"id"`
	AdvertisementId string `json:"advertisement_id"`
	Url             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
	CreatedAt       string `json:"created_at"`
}
