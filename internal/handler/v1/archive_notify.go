package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerArchiveNotify(router *gin.RouterGroup) {
	route := router.Group("/archive_notify")
	// route.POST("", h.CreateNotify)
	route.POST("/find", h.FindArchiveNotifyPopulate)
	route.DELETE("/:id", h.DeleteArchiveNotify)
}

// func (h *HandlerV1) CreateNotify(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	// userID, err := middleware.GetUID(c)
// 	// if err != nil {
// 	// 	// c.AbortWithError(http.StatusUnauthorized, err)
// 	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 	// 	return
// 	// }

// 	var input *domain.NotifyInput
// 	if er := c.Bind(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	notify, err := h.CreateOrExistNotify(c, input) //h.services.Notify.CreateNotify(userID, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, notify)
// }

// @Summary Find Notifys by params
// @Security ApiKeyAuth
// @Tags Notify
// @Description Input params for search Notifys
// @ModuleID Notify
// @Accept  json
// @Produce  json
// @Param input query NotifyInput true "params for search Notify"
// @Success 200 {object} []domain.Notify
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Notify [get].
func (h *HandlerV1) FindArchiveNotifyPopulate(c *gin.Context) {
	appG := app.Gin{C: c}
	// params, err := utils.GetParamsFromRequest(c, domain.NotifyInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.ArchiveNotifyFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Notifys, err := h.Services.ArchiveNotify.FindArchiveNotifyPopulate(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Notifys)
}

// func (h *HandlerV1) CreateOrExistNotify(c *gin.Context, input *domain.NotifyInput) (*domain.Notify, error) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return nil, err
// 	}
// 	var result *domain.Notify

// 	// upload images.
// 	var imageInput = &domain.MessageImage{}
// 	imageInput.Service = "notify"
// 	imageInput.ServiceID = primitive.NilObjectID.Hex()
// 	imageInput.UserID = userID

// 	paths, err := utils.UploadResizeMultipleFileForMessage(c, imageInput, "images", &h.imageConfig)
// 	if err != nil {
// 		appG.ResponseError(http.StatusInternalServerError, err, nil)
// 	}

// 	resultImages := []string{}
// 	for i := range paths {
// 		imageInput.Path = paths[i].Path
// 		imageInput.Ext = paths[i].Ext
// 		// imageInput.Service= "message"
// 		// image, err := h.Services.MessageImage.CreateMessageImage(userID, imageInput)
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return result, err
// 		// }
// 		// imageInput.URL =
// 		resultImages = append(resultImages, fmt.Sprintf("%s/%s/%s/%s%s", imageInput.UserID, imageInput.Service, imageInput.ServiceID, imageInput.Path, imageInput.Ext))
// 	}

// 	input.Images = resultImages

// 	// create notify.
// 	result, err = h.Services.Notify.CreateNotify(userID, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}

// 	return result, nil
// }

func (h *HandlerV1) DeleteArchiveNotify(c *gin.Context) {
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

	user, err := h.Services.ArchiveNotify.DeleteArchiveNotify(id, userID) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
