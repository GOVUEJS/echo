package model

type RedisSession struct {
	Email        *string
	Ip           *string
	AccessToken  *string
	RefreshToken *string
}
