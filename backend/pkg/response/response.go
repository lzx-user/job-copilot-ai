package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Body{
		Code: "OK",
		Data: data,
	})
}

func Error(ctx *gin.Context, httpStatus int, code string, message string) {
	ctx.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
