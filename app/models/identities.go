package models

type Identity struct {
	Network         string `json:"network"            bson:"network"`
	Id              string `json:"id"                 bson:"id"`
	Username        string `json:"username"           bson:"username"`
	ProfileURL      string `json:"profile_url"        bson:"profile_url"`
	ProfileImageURL string `json:"profile_image_url"  bson:"profile_image_url"`
	Status          string `json:"status"             bson:"status"`
}

type IdentityResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Identity `json:"data,omitempty"`
}
