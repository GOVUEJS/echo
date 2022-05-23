package util

import (
	"github.com/golang-jwt/jwt"
	"myapp/config"
	"myapp/consts"
	"myapp/model"
	"time"
)

func GetAccessRefreshToken(email, sessionId *string) (accessToken, refreshToken *string, err error) {

	// Set custom accessTokenClaims
	accessTokenClaims := &model.JwtCustomClaims{
		SessionId: *sessionId,
		Email:     *email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.SessionDuration / 10).Unix(),
		},
	}

	// Create token with accessTokenClaims
	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// Generate encoded token and send it as response.
	accessTokenString, err := accessTokenJWT.SignedString(config.Config.Jwt.Key)
	if err != nil {
		return nil, nil, err
	}

	// Set custom refreshTokenClaims
	refreshTokenClaims := &model.JwtCustomClaims{
		SessionId: *sessionId,
		Email:     *email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.SessionDuration).Unix(),
		},
	}

	// Create token with refreshTokenClaims
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Generate encoded token and send it as response.
	refreshTokenString, err := refreshTokenJWT.SignedString(config.Config.Jwt.Key)
	if err != nil {
		return nil, nil, err
	}

	accessToken = &accessTokenString
	refreshToken = &refreshTokenString
	return
}
