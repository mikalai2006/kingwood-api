package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerArchiveTask(router *gin.RouterGroup) {
	route := router.Group("/archive_task")
	route.GET("", h.FindArchiveTask)
	route.DELETE("/:id", h.SetUserFromRequest, h.DeleteArchiveTask)
}

func (h *HandlerV1) FindArchiveTask(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.ArchiveTaskFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	result, err := h.Services.ArchiveTask.FindArchiveTask(*input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *HandlerV1) DeleteArchiveTask(c *gin.Context) {
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
		return
	}

	user, err := h.Services.ArchiveTask.DeleteArchiveTask(id, userID)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
