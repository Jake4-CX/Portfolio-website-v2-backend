package structs

import ()

type LoginResponseModel struct {
	User  Users       `json:"user"`
	Token TokensModel `json:"token"`
}

type TokensModel struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
