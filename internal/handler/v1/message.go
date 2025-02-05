package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerMessage(router *gin.RouterGroup) {
	message := router.Group("/message")
	message.POST("", h.CreateMessage)
	message.POST("/list", h.CreateListMessage)
	message.POST("/find", h.FindMessage)
	message.PATCH("/:id", h.UpdateMessage)
	message.DELETE("/:id", h.DeleteMessage)
	message.GET("/groups", h.GetGroupByUser)
}

func (h *HandlerV1) CreateMessage(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.MessageInput
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// node, err := h.services.Message.CreateMessage(userID, input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	node, err := h.CreateOrExistMessage(c, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateListMessage(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input []*domain.MessageInput
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Message
	for i := range input {
		Message, err := h.CreateOrExistMessage(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Message)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find Messages by params
// @Security ApiKeyAuth
// @Tags Message
// @Description Input params for search Messages
// @ModuleID Message
// @Accept  json
// @Produce  json
// @Param input query Message true "params for search Message"
// @Success 200 {object} []domain.Message
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/node_audit [get].
func (h *HandlerV1) FindMessage(c *gin.Context) {
	appG := app.Gin{C: c}

	// authData, err := middleware.GetAuthFromCtx(c)
	// fmt.Println("auth ", authData.Roles)

	// params, err := utils.GetParamsFromRequest(c, domain.Message{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.MessageFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(params)
	Nodes, err := h.Services.Message.FindMessage(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodes)
}

func (h *HandlerV1) UpdateMessage(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// var input model.TagInput
	// data, err := utils.BindAndValid(c, &input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[model.Node](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// fmt.Println(data)
	var input *domain.MessageInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Message.UpdateMessage(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteMessage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	// implementation roles for user.
	roles, err := middleware.GetRoles(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if !utils.Contains(roles, "admin") {
		appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
		return
	}

	node, err := h.Services.Message.DeleteMessage(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateOrExistMessage(c *gin.Context, input *domain.MessageInput) (*domain.Message, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	// nodeIDPrimitive, err := primitive.ObjectIDFromHex(string(input.NodeID))
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return nil, err
	// }
	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return nil, err
	// }

	var result *domain.Message

	// check exist order.
	orderID, _ := primitive.ObjectIDFromHex(input.OrderID)
	fmt.Println(orderID, input)
	existRoom, err := h.Services.Order.FindOrder(&domain.OrderFilter{ID: []string{orderID.Hex()}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	if len(existRoom.Data) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("not found room for message"), nil)
		return result, nil
	}

	// input.UserProductID = existNode.Data[0].UserID

	// // check exist message
	// existMessage, err := h.services.Message.FindMessage(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"product_id", input.ProductID},
	// 		{"userId", userIDPrimitive},
	// 	},
	// 	Options: domain.Options{
	// 		Limit: 1,
	// 	},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return nil, err
	// }
	// if len(existMessage.Data) > 0 {
	// 	// //appG.ResponseError(http.StatusBadRequest, errors.New("existSameNode"), nil)
	// 	// update node audit.
	// 	id := &existMessage.Data[0].ID
	// 	result, err = h.services.Message.UpdateMessage(id.Hex(), userID, input)
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return result, err
	// 	}

	// 	return result, nil
	// } else {
	// }

	// upload images.
	var imageInput = &domain.MessageImage{}
	imageInput.Service = "message"
	imageInput.ServiceID = input.OrderID
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

	// create message.
	result, err = h.Services.Message.CreateMessage(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	return result, nil
}

func (h *HandlerV1) GetGroupByUser(c *gin.Context) {
	appG := app.Gin{C: c}

	// id := c.Param("id")
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	groups, err := h.Services.Message.GetGroupForUser(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, groups)
}
