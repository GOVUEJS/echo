package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"myapp/config"
	"myapp/database/rdb/postgres"
	"myapp/database/session"
	"myapp/logger"
	"myapp/service"
	v1 "myapp/service/v1"
)

func newEcho() *echo.Echo {
	logger.InitLogger()

	err := config.InitConfig("../config/echo-dev.toml")
	if err != nil {
		log.Fatalf("Decode toml error: %s", err)
		panic(err)
	}

	err = session.InitRedis()
	if err != nil {
		log.Fatalf("InitRedis error: %s", err)
		panic(err)
	}

	err = postgres.InitRDB()
	if err != nil {
		log.Fatalf("InitPostgres error: %s", err)
		panic(err)
	}

	service.InitService()

	return echo.New()
}

func TestPostLogin(t *testing.T) {
	// Setup
	e := newEcho()

	bodyJSON := `{"email":"test@test.com","pw":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(bodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, v1.PostLogin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
