package models

type Identity struct {
	Network         string `json:"network"`
	Id              string `json:"id"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
	Status          string `json:"status"`
}

type IdentityResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Identity `json:"data,omitempty"`
}
