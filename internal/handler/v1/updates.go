package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) registerUpdates(router *gin.RouterGroup) {
	route := router.Group("/updates")
	route.GET("/manifest", h.Manifest)
}

func (h *HandlerV1) Manifest(c *gin.Context) {
	// appG := app.Gin{C: c}

	fmt.Println("Manifest: headers - ", c.Request.Header)

	response := map[string]interface{}{}
	c.JSON(http.StatusOK, response)
}
