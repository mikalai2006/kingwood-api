package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerTimer(router *gin.RouterGroup) {
	route := router.Group("/timer")
	route.POST("/populate", h.FindTimerPopulate)
	route.DELETE("/:id", h.SetUserFromRequest, h.DeleteTimer)
	// route.POST("/recovery", h.RecoveryTimers)
}

func (h *HandlerV1) FindTimerPopulate(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.TimerSheduleFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Orders, err := h.Services.Timer.FindTimerPopulate(*input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Orders)
}

func (h *HandlerV1) DeleteTimer(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

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

	user, err := h.Services.Timer.DeleteTimer(id, userID)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// func (h *HandlerV1) RecoveryTimers(c *gin.Context) {
// 	// appG := app.Gin{C: c}
// 	// userID, err := middleware.GetUID(c)
// 	// if err != nil {
// 	// 	// c.AbortWithError(http.StatusUnauthorized, err)
// 	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 	// 	return
// 	// }
// 	fmt.Println("Recovery timers")

// 	// userID, err := middleware.GetUID(c)
// 	// if err != nil {
// 	// 	// c.AbortWithError(http.StatusUnauthorized, err)
// 	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 	// 	return
// 	// }

// 	// user, err := h.Services.Timer.DeleteTimer(id, userID)
// 	// if err != nil {
// 	// 	// c.AbortWithError(http.StatusBadRequest, err)
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return
// 	// }

// 	c.JSON(http.StatusOK, nil)
// }
