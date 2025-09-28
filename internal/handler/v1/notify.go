package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerNotify(router *gin.RouterGroup) {
	route := router.Group("/notify")
	route.POST("", h.CreateNotify)
	route.POST("/list", h.CreateNotifyList)
	route.POST("/populate", h.FindNotifyPopulate)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateNotify)
	route.DELETE("/:id", h.DeleteNotify)
	route.POST("/remove_list", h.RemoveNotifyList)
	route.POST("/clear", h.ClearNotify)
}

func (h *HandlerV1) CreateNotify(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.NotifyInput
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	notify, err := h.CreateOrExistNotify(c, input) //h.services.Notify.CreateNotify(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, notify)
}

func (h *HandlerV1) CreateNotifyList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.NotifyInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Notify
	for i := range input {
		Notify, err := h.CreateOrExistNotify(c, input[i]) //h.services.Notify.CreateNotify(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Notify)
	}

	c.JSON(http.StatusOK, result)
}

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
func (h *HandlerV1) FindNotifyPopulate(c *gin.Context) {
	appG := app.Gin{C: c}
	// params, err := utils.GetParamsFromRequest(c, domain.NotifyInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.NotifyFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Notifys, err := h.Services.Notify.FindNotifyPopulate(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Notifys)
}

func (h *HandlerV1) UpdateNotify(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// // var input domain.PostInput
	// // data, err := utils.BindAndValid(c, &input)
	// // if err != nil {
	// // 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// // 	return
	// // }
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[domain.Post](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// // fmt.Println(data)

	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[domain.NotifyInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Notify.UpdateNotify(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) CreateOrExistNotify(c *gin.Context, input *domain.NotifyInput) (*domain.Notify, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.Notify

	// upload images.
	var imageInput = &domain.MessageImage{}
	imageInput.Service = "notify"
	imageInput.ServiceID = primitive.NilObjectID.Hex()
	imageInput.UserID = userID

	paths, err := utils.UploadResizeMultipleFileForMessage(c, imageInput, "images", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}

	resultImages := []string{}
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		// imageInput.Service= "message"
		// image, err := h.Services.MessageImage.CreateMessageImage(userID, imageInput)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return result, err
		// }
		// imageInput.URL =
		resultImages = append(resultImages, fmt.Sprintf("%s/%s/%s/%s%s", imageInput.UserID, imageInput.Service, imageInput.ServiceID, imageInput.Path, imageInput.Ext))
	}

	input.Images = resultImages

	// create notify.
	result, err = h.Services.Notify.CreateNotify(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	return result, nil
}

func (h *HandlerV1) DeleteNotify(c *gin.Context) {
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

	user, err := h.Services.Notify.DeleteNotify(id, userID, true) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) RemoveNotifyList(c *gin.Context) {
	appG := app.Gin{C: c}
	var result *[]domain.Notify

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	user, err := h.Services.User.GetUser(userID)

	if user.RoleObject.Code == "systemrole" {
		var input *domain.NotifyListQuery
		if er := c.BindJSON(&input); er != nil {
			appG.ResponseError(http.StatusBadRequest, er, nil)
			return
		}

		result, err = h.Services.Notify.DeleteNotifyList(*input)
		if err != nil {
			// c.AbortWithError(http.StatusUnauthorized, err)
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
	} else {
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *HandlerV1) ClearNotify(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	user, err := h.Services.User.GetUser(userID)

	if user.RoleObject.Code == "systemrole" {
		err = h.Services.Notify.ClearNotify(userID)
	} else {
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
