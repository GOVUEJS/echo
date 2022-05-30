package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"myapp/database/rdb/postgres"
	"myapp/model"
	"myapp/service/v1"
	"myapp/test"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestGetArticleList(t *testing.T) {
	type args struct {
		Page int
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "page=1 - 200",
			args: args{
				Page: 1,
			},
			wantResult: http.StatusOK,
		},
		{
			name: "page=-1 - 400",
			args: args{
				Page: -1,
			},
			wantResult: http.StatusBadRequest,
		},
		{
			name: "page=10 - 200",
			args: args{
				Page: 10,
			},
			wantResult: http.StatusOK,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/articles"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := make(url.Values)
			query.Set("page", strconv.Itoa(tt.args.Page))

			marshal, _ := json.Marshal(&tt.args)

			req := httptest.NewRequest(http.MethodGet, target+"/?"+query.Encode(), strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			if assert.NoError(t, v1.GetArticleList(c)); rec.Code != tt.wantResult {
				t.Errorf("GetArticleList() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestGetArticle(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "Id=1 - 404",
			args: args{
				id: "1",
			},
			wantResult: http.StatusNotFound,
		},
		{
			name: "Id=-1 - 400",
			args: args{
				id: "-1",
			},
			wantResult: http.StatusBadRequest,
		},
		{
			name: "Id=8 - 200",
			args: args{
				id: "8",
			},
			wantResult: http.StatusOK,
		},
		{
			name: "Id=q - 400",
			args: args{
				id: "q",
			},
			wantResult: http.StatusBadRequest,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/article/:id"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodGet, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.args.id)

			// Assertions
			if assert.NoError(t, v1.GetArticle(c)); rec.Code != tt.wantResult {
				t.Errorf("GetArticle() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestPostArticle(t *testing.T) {
	type args struct {
		Title   string
		Content string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "",
			args: args{
				Title:   "test",
				Content: "test",
			},
			wantResult: http.StatusCreated,
		},
		{
			name: "",
			args: args{
				Title:   "",
				Content: "test",
			},
			wantResult: http.StatusBadRequest,
		},
		{
			name: "",
			args: args{
				Title:   "test",
				Content: "",
			},
			wantResult: http.StatusBadRequest,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/articles"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("email", "test@test.com")

			// Assertions
			if assert.NoError(t, v1.PostArticle(c)); rec.Code != tt.wantResult {
				t.Errorf("PostArticle() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestPutArticle(t *testing.T) {
	type args struct {
		id      string
		Title   string
		Content string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "",
			args: args{
				id:      "7",
				Title:   "test",
				Content: "test",
			},
			wantResult: http.StatusNotFound,
		},
		{
			name: "",
			args: args{
				id:      "78",
				Title:   "",
				Content: "test",
			},
			wantResult: http.StatusBadRequest,
		},
		{
			name: "",
			args: args{
				id:      "78",
				Title:   "test",
				Content: "",
			},
			wantResult: http.StatusBadRequest,
		},
		{
			name: "",
			args: args{
				id:      "78",
				Title:   "test",
				Content: uuid.New().String(),
			},
			wantResult: http.StatusOK,
		},
		{
			name: "",
			args: args{
				id:      "73",
				Title:   "test",
				Content: uuid.New().String(),
			},
			wantResult: http.StatusUnauthorized,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/articles/:id"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodPut, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("email", "test@test.com")
			c.SetParamNames("id")
			c.SetParamValues(tt.args.id)

			// Assertions
			if assert.NoError(t, v1.PutArticle(c)); rec.Code != tt.wantResult {
				t.Errorf("PutArticle() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestDeleteArticle(t *testing.T) {

	type args struct {
		Id string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "removeOK",
			args: args{
				Id: "",
			},
			wantResult: http.StatusNoContent,
		},
	}

	e := test.NewEchoForTest()
	target := "/api/v1/articles/:id"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "removeOK" {
				postgres.GetRDB().Model(model.Article{}).Select("id").Where("writer = ?", "test@test.com").Order("id desc").First(&tt.args.Id)
			}

			marshal, _ := json.Marshal(&tt.args)
			req := httptest.NewRequest(http.MethodDelete, target, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("email", "test@test.com")
			c.SetParamNames("id")
			c.SetParamValues(tt.args.Id)

			// Assertions
			if assert.NoError(t, v1.DeleteArticle(c)); rec.Code != tt.wantResult {
				t.Errorf("DeleteArticle() gotResult = %v, want = %v, msg = %v", rec.Code, tt.wantResult, rec.Body.String())
			}
		})
	}
}
