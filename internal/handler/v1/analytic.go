package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) RegisterAnalytic(router *gin.RouterGroup) {
	route := router.Group("/analytic")
	route.GET("/get", h.getAnalytic)
}

func (h *HandlerV1) getAnalytic(c *gin.Context) {
	appG := app.Gin{C: c}

	document, err := h.Services.Analytic.GetAnalytic()
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
