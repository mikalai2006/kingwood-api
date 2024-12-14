package domain

import "time"

type Response[D any] struct {
	Total int `json:"total" bson:"total"`
	Limit int `json:"limit" bson:"limit"`
	Skip  int `json:"skip" bson:"skip"`
	Data  []D `json:"data" bson:"data"`
}

type ResponseTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in" bson:"expires_in"`
	ExpiresInR   int64  `json:"expires_in_r" bson:"expires_in_r"`
}

type GeneralFieldDB struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" form:"updatedAt"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseNominatim struct {
	PlaceId     int                    `json:"place_id" bson:"place_id"`
	OsmType     string                 `json:"osm_type" bson:"osm_type"`
	OsmID       int                    `json:"osm_id" bson:"osm_id"`
	Lat         string                 `json:"lat" bson:"lat"`
	Lon         string                 `json:"lon" bson:"lon"`
	Class       string                 `json:"class" bson:"class"`
	Type        string                 `json:"type" bson:"type"`
	Lang        string                 `json:"lang" bson:"lang"`
	AddressType string                 `json:"addresstype" bson:"addresstype"`
	Name        string                 `json:"name" bson:"name"`
	DisplayName string                 `json:"display_name" bson:"display_name"`
	Address     map[string]interface{} `json:"address" bson:"address"`
}

// type ResponseNominatimAddress struct {
// 	Highway      string `json:"highway"`
// 	Road         string `json:"road"`
// 	Village      string `json:"village"`
// 	Municipality string `json:"municipality"`
// 	County       string `json:"county"`
// 	State        string `json:"state"`
// 	Postcode     string `json:"postcode"`
// 	Country      string `json:"country"`
// 	CountryCode  string `json:"country_code"`
// }
