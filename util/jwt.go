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

func CheckToken(tokens *model.Tokens) (accessTokenClaims, refreshTokenClaims jwt.MapClaims, err error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return config.Config.Jwt.Key, nil
	}

	accessToken, err := jwt.Parse(*tokens.AccessToken, keyFunc)
	if err != nil {
		return nil, nil, err
	}
	if !accessToken.Valid {
		return nil, nil, errors.New("invalid tokens")
	}
	accessTokenClaims, _ = accessToken.Claims.(jwt.MapClaims)

	refreshToken, err := jwt.Parse(*tokens.RefreshToken, keyFunc)
	if err != nil {
		return nil, nil, err
	}
	if !accessToken.Valid {
		return nil, nil, errors.New("invalid tokens")
	}
	refreshTokenClaims, _ = refreshToken.Claims.(jwt.MapClaims)

	if accessTokenClaims["email"] != refreshTokenClaims["email"] || accessTokenClaims["sessionId"] != refreshTokenClaims["sessionId"] {
		return nil, nil, errors.New("invalid tokens")
	}

	return accessTokenClaims, refreshTokenClaims, nil
}
