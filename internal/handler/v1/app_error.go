package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerAppError(router *gin.RouterGroup) {
	route := router.Group("/error_app")
	route.POST("", h.CreateAppError)
	route.POST("/populate", h.FindAppError)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateAppError)
	route.DELETE("/:id", h.DeleteAppError)
}

func (h *HandlerV1) CreateAppError(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.AppError
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	AppError, err := h.CreateOrExistAppError(c, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, AppError)
}

func (h *HandlerV1) CreateAppErrorList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.AppError
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.AppError
	for i := range input {
		AppError, err := h.CreateOrExistAppError(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, AppError)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find AppErrors by params
// @Security ApiKeyAuth
// @Tags AppError
// @Description Input params for search AppErrors
// @ModuleID AppError
// @Accept  json
// @Produce  json
// @Param input query AppErrorInput true "params for search AppError"
// @Success 200 {object} []domain.AppError
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/AppError [get].
func (h *HandlerV1) FindAppError(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.AppErrorFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	AppErrors, err := h.Services.AppError.FindAppError(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, AppErrors)
}

func (h *HandlerV1) GetAppErrorByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateAppError(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	var a map[string]json.RawMessage
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[domain.AppErrorInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.AppError.UpdateAppError(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteAppError(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	user, err := h.Services.AppError.DeleteAppError(id, userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) CreateOrExistAppError(c *gin.Context, input *domain.AppError) (*domain.AppError, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.AppError

	result, err = h.Services.AppError.CreateAppError(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}
