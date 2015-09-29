package models

import ()

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Not common but we need them to query the tickets API
type CredentialRequestData struct {
	Server  string `json:"server"`
	Tickets int    `json:"tickets"`
}

type CredentialResponseData struct {
	Tokens  string `json:"tokens"`
	Tickets int    `json:"tickets"`
}

type CredentialResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    CredentialResponseData `json:"data,omitempty"`
}
