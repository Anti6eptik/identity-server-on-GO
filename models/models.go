package models


type User struct{
	UserName string
	Password string
}

type Tokens struct{
	AccessToken string
	RefreshToken string
}

type TokenModel struct{
	Token string `json:"access_token"`
}