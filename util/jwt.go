package util

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"myapp/config"
	"myapp/consts"
	"myapp/model"
	"time"
)

func GetAccessRefreshToken(email, sessionId *string) (accessToken, refreshToken *string, err error) {
	return GetAccessRefreshTokenWithDuration(email, sessionId, consts.SessionDuration)
}

func GetAccessRefreshTokenWithDuration(email, sessionId *string, duration time.Duration) (accessToken, refreshToken *string, err error) {

	// Set custom accessTokenClaims
	accessTokenClaims := &model.JwtCustomClaims{
		SessionId: *sessionId,
		Email:     *email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration / 10).Unix(),
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
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	// Create token with refreshTokenClaims
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Generate encoded token and send it as response.
	refreshTokenString, err := refreshTokenJWT.SignedString(config.Config.Jwt.Key)
	if err != nil {
		return nil, nil, err
	}

	return &accessTokenString, &refreshTokenString, nil
}

func CheckRefreshToken(tokens *model.Tokens) (accessTokenClaims, refreshTokenClaims jwt.MapClaims, err error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return config.Config.Jwt.Key, nil
	}

	accessToken, err := jwt.Parse(*tokens.AccessToken, keyFunc)
	accessTokenClaims, _ = accessToken.Claims.(jwt.MapClaims)

	refreshToken, err := jwt.Parse(*tokens.RefreshToken, keyFunc)
	if err != nil {
		return nil, nil, err
	}
	if !refreshToken.Valid {
		return nil, nil, errors.New("invalid tokens")
	}
	refreshTokenClaims, _ = refreshToken.Claims.(jwt.MapClaims)

	if accessTokenClaims["email"] != refreshTokenClaims["email"] || accessTokenClaims["sessionId"] != refreshTokenClaims["sessionId"] {
		return nil, nil, errors.New("invalid tokens")
	}

	return accessTokenClaims, refreshTokenClaims, nil
}
