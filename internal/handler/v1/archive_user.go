package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) RegisterArchiveUser(router *gin.RouterGroup) {
	route := router.Group("/archive_user")
	// route.POST("", h.CreateUser)
	route.POST("/find", h.FindArchiveUser)
	route.DELETE("/:id", h.DeleteArchiveUser)
}

func (h *HandlerV1) FindArchiveUser(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.UserInput{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.ArchiveUserFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	users, err := h.Services.ArchiveUser.FindArchiveUser(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

// func (h *HandlerV1) CreateUser(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		appG.ResponseError(http.StatusUnauthorized, err, nil)
// 		return
// 	}

// 	var input *domain.User
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	user, er := h.Services.User.CreateUser(userID, input)
// 	if er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Delete user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [delete].
func (h *HandlerV1) DeleteArchiveUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	user, err := h.Services.ArchiveUser.DeleteArchiveUser(id, userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
