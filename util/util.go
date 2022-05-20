package util

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myapp/consts"
	"myapp/model"
	"net/http"
)

func Response(c echo.Context, status int, message string, data interface{}) error {
	response := make(map[string]interface{})

	m2, _ := json.Marshal(data)
	if err := json.Unmarshal(m2, &response); err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Error")
	}

	apiResponse := model.ApiResponse{}
	if message != "" {
		apiResponse.Message = &message
	}
	m1, _ := json.Marshal(apiResponse)
	if err := json.Unmarshal(m1, &response); err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Error")
	}

	return c.JSON(status, response)
}

func GetPagination(sql *gorm.DB, page int) (current, total int64) {
	var totalCount int64
	sql.Count(&totalCount)
	total = GetTotalPage(totalCount)
	current = int64(page)
	return
}

func GetTotalPage(totalCount int64) (totalPage int64) {
	totalPage = totalCount / consts.PageSize
	if totalCount%consts.PageSize != 0 {
		totalPage++
	}
	return
}

func GetPageOffsetLimit(page int) (offset, limit int) {
	return (page - 1) * consts.PageSize, consts.PageSize
}
