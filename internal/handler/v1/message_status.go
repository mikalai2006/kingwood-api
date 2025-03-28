package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerMessageStatus(router *gin.RouterGroup) {
	MessageStatus := router.Group("/message_status")
	MessageStatus.POST("", h.CreateMessageStatus)
	MessageStatus.POST("/find", h.FindMessageStatus)
	MessageStatus.PATCH("/:id", h.UpdateMessageStatus)
	MessageStatus.DELETE("/:id", h.DeleteMessageStatus)
}

func (h *HandlerV1) CreateMessageStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.MessageStatus
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// node, err := h.services.Message.CreateMessage(userID, input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	node, err := h.CreateOrExistMessageStatus(c, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
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
func (h *HandlerV1) FindMessageStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	// authData, err := middleware.GetAuthFromCtx(c)
	// fmt.Println("auth ", authData.Roles)

	// params, err := utils.GetParamsFromRequest(c, domain.Message{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.MessageStatusFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(params)
	Nodes, err := h.Services.MessageStatus.FindMessageStatus(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodes)
}

func (h *HandlerV1) UpdateMessageStatus(c *gin.Context) {
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
	var input *domain.MessageStatus
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.MessageStatus.UpdateMessageStatus(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteMessageStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	// // implementation roles for user.
	// roles, err := middleware.GetRoles(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// if !utils.Contains(roles, "admin") {
	// 	appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
	// 	return
	// }

	node, err := h.Services.MessageStatus.DeleteMessageStatus(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateOrExistMessageStatus(c *gin.Context, input *domain.MessageStatus) (*domain.MessageStatus, error) {
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

	var result *domain.MessageStatus

	// check exist product.
	// domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"_id", input.ProductID},
	// 	},
	// 	Options: domain.Options{
	// 		Limit: 1,
	// 	},
	// }
	messageID := input.MessageID.Hex()
	existMessage, err := h.Services.Message.FindMessage(&domain.MessageFilter{ID: messageID})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	if len(existMessage.Data) == 0 {
		//appG.ResponseError(http.StatusBadRequest, errors.New("not found node"), nil)
		return result, nil
	}

	input.MessageID = existMessage.Data[0].ID

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
	// create message.
	result, err = h.Services.MessageStatus.CreateMessageStatus(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	return result, nil
}
