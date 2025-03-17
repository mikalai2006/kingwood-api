package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerArchiveTaskWorker(router *gin.RouterGroup) {
	route := router.Group("/archive_task_worker")
	route.POST("/find", h.FindArchiveTaskWorker)
	route.DELETE("/:id", h.DeleteArchiveTaskWorker)
}

func (h *HandlerV1) FindArchiveTaskWorker(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.ArchiveTaskWorkerFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	ArchiveTaskWorkers, err := h.Services.ArchiveTaskWorker.FindArchiveTaskWorker(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, ArchiveTaskWorkers)
}

func (h *HandlerV1) DeleteArchiveTaskWorker(c *gin.Context) {
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

	user, err := h.Services.ArchiveTaskWorker.DeleteArchiveTaskWorker(id, userID) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
