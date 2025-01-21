package httputil

import "github.com/gin-gonic/gin"

func NewResponse(ctx *gin.Context, code int, status, message string, data any) {
	response := HTTPResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if data != nil {
		ctx.JSON(code, data)
	} else {
		ctx.JSON(code, response)
	}

	ctx.Abort()
}

type HTTPResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
