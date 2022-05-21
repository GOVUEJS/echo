package model

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
