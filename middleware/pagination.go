package middleware

import (
	"gin-photo-storage/constant"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var InvalidPageNoError = errors.New("page no can not be negative")

// A wrapper function which returns the pagination middleware.
func GetPaginationMiddleware() func(*gin.Context) {
	return func(context *gin.Context) {
		responseCode := constant.PAGINATION_SUCCESS
		pageNo := context.Query("page")
		if pageNo == "" {
			responseCode = constant.INVALID_PARAMS
		} else {
			pageOffset, err := GetPagination(pageNo)
			if err != nil {
				responseCode = constant.INVALID_PARAMS
			} else {
				context.Set("offset", pageOffset)
			}
		}

		if responseCode == constant.INVALID_PARAMS {
			data := make(map[string]string)
			data["page"] = pageNo
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.INVALID_PARAMS,
				"data": data,
				"msg":  constant.GetMessage(constant.INVALID_PARAMS),
			})
			context.Abort()
		}

		context.Next()
	}
}

// Pagination function which calculates the offset given the page number.
func GetPagination(pageNo string) (int, error) {
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return 0, err
	}
	if pageNoInt < 0 {
		return 0, InvalidPageNoError
	}
	return pageNoInt * constant.PAGE_SIZE, nil
}