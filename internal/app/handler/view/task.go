package view

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Title struct {
	Title string `json:"title"`
}

func Bind[R any](c *gin.Context) (R, error) {
	var request R
	err := c.ShouldBindJSON(&request)
	return request, err
}

type PaginationQuery struct {
	Page  string `query:"page"`
	Limit string `query:"limit"`
	Order string `query:"order"`
}

func (p PaginationQuery) GetPage() int {
	page, err := strconv.Atoi(p.Page)
	if err == nil && page > 0 {
		return page
	}
	return 1
}

func (p PaginationQuery) GetLimit() int {
	limit, err := strconv.Atoi(p.Limit)
	if err != nil || limit < 1 {
		return 200
	} else if limit > 5000 {
		return 5000
	}
	return limit
}

func (p PaginationQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
