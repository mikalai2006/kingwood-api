package app

import "github.com/gin-gonic/gin"

type Gin struct {
	C *gin.Context
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (g *Gin) ResponseError(httpCode int, err error, data interface{}) {
	g.C.JSON(httpCode, ErrorResponse{
		Code:    httpCode,
		Message: err.Error(),
		Data:    data,
	})
	g.C.Abort()
	// or g.C.AbortWithError(httpCode, err)
}
