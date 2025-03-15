package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerArchiveOrder(router *gin.RouterGroup) {
	route := router.Group("/archive_order")
	route.POST("/find", h.FindArchiveOrder)
	route.DELETE("/:id", h.DeleteArchiveOrder)
}

// func (h *HandlerV1) CreateArchiveOrder(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil || userID == "" {
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return
// 	}

// 	var input *domain.ArchiveOrder
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	result, err := h.Services.ArchiveOrder.CreateArchiveOrder(userID, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, result)
// }

func (h *HandlerV1) FindArchiveOrder(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.ArchiveOrderFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	archiveOrders, err := h.Services.ArchiveOrder.FindArchiveOrder(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, archiveOrders)
}

func (h *HandlerV1) DeleteArchiveOrder(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	}

	user, err := h.Services.ArchiveOrder.DeleteArchiveOrder(id, userID) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
