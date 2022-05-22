package util

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/random"
	"myapp/model"
	"time"
)

var (
	JwtKey []byte
)

func init() {

	JwtKey = []byte(random.String(32))
}
func GetAccessRefreshToken(email, sessionId *string) (accessToken, refreshToken *string, err error) {

	// Set custom accessTokenClaims
	accessTokenClaims := &model.JwtCustomClaims{
		SessionId: *sessionId,
		Email:     *email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with accessTokenClaims
	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// Generate encoded token and send it as response.
	accessTokenString, err := accessTokenJWT.SignedString(JwtKey)
	if err != nil {
		return nil, nil, err
	}

	// Set custom refreshTokenClaims
	refreshTokenClaims := &model.JwtCustomClaims{
		SessionId: *sessionId,
		Email:     *email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// Create token with refreshTokenClaims
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Generate encoded token and send it as response.
	refreshTokenString, err := refreshTokenJWT.SignedString(JwtKey)
	if err != nil {
		return nil, nil, err
	}

	accessToken = &accessTokenString
	refreshToken = &refreshTokenString
	return
}
