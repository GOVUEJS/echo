package v1

import (
	"encoding/json"
	"myapp/test"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"myapp/service/v1"
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
