package test

import (
	"encoding/json"
	"flag"
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

var (
	filePath = flag.String("configFilePath", "", "configFilePath")
)

func newEcho() *echo.Echo {
	logger.InitLogger()

	err := config.InitConfig(*filePath)
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

func TestPostSignUp(t *testing.T) {
	type args struct {
		Email string
		Pw    string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "test@test.com - 400",
			args: args{
				Email: "test@test.com",
				Pw:    "test",
			},
			wantResult: http.StatusBadRequest,
		},
	}

	e := newEcho()
	target := "/api/v1/signup"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, v1.PostSignUp(c)); rec.Code != tt.wantResult {
				t.Errorf("PostSignUp() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestPostLogin(t *testing.T) {
	type args struct {
		Email string
		Pw    string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "test@test.com - 200",
			args: args{
				Email: "test@test.com",
				Pw:    "test",
			},
			wantResult: http.StatusOK,
		},
		{
			name: "test@test.com - 400",
			args: args{
				Email: "test@test.com",
				Pw:    "test1",
			},
			wantResult: http.StatusBadRequest,
		},
	}

	e := newEcho()
	target := "/api/v1/login"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, v1.PostLogin(c)); rec.Code != tt.wantResult {
				t.Errorf("PostLogin() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}
