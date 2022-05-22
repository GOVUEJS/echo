package model

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	SessionId string `json:"sessionId"`
	Email     string `json:"email"`
	jwt.StandardClaims
}
