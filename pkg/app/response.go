package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Gin struct {
	C *gin.Context
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (g *Gin) ResponseError(httpCode int, err error, data interface{}) {

	// if err != nil && err != mongo.ErrNoDocuments {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	errorText := ""
	if err != nil {
		errorText = err.Error()
		if err == mongo.ErrNoDocuments {
			httpCode = http.StatusUnauthorized
			errorText = domain.ErrNotItemMongo.Error()
		}
	}

	g.C.JSON(httpCode, ErrorResponse{
		Code:    httpCode,
		Message: errorText,
		Data:    data,
	})
	g.C.Abort()
	// or g.C.AbortWithError(httpCode, err)
}
