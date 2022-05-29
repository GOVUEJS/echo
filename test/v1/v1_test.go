package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"myapp/model"
	"myapp/service/v1"
	"myapp/test"
	"myapp/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

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

	e := test.NewEchoForTest()
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

	e := test.NewEchoForTest()
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

func TestGetLogout(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name:       testing.CoverMode(),
			args:       args{},
			wantResult: http.StatusOK,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/logout"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, v1.GetLogout(c)); rec.Code != tt.wantResult {
				t.Errorf("GetLogout() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestPostRefreshToken(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name:       testing.CoverMode(),
			args:       args{},
			wantResult: http.StatusOK,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/token/refresh"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email := "test@test.com"
			sessionId := uuid.New().String()
			accessToken, refreshToken, _ := util.GetAccessRefreshTokenWithDuration(&email, &sessionId, time.Second*10)
			tokens := &model.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}

			marshal, _ := json.Marshal(tokens)
			req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, v1.PostRefreshToken(c)); rec.Code != tt.wantResult {
				t.Errorf("PostRefreshToken()() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}
