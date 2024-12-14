package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerMessageRoom(router *gin.RouterGroup) {
	messageRoom := router.Group("/message_room")
	messageRoom.POST("", h.CreateMessageRoom)
	messageRoom.POST("/find", h.FindMessageRoom)
	messageRoom.PATCH("/:id", h.UpdateMessageRoom)
	messageRoom.DELETE("/:id", h.DeleteMessageRoom)
}

func (h *HandlerV1) CreateMessageRoom(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *model.MessageRoom
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// node, err := h.services.Message.CreateMessage(userID, input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	node, err := h.CreateOrExistMessageRoom(c, input)
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
func (h *HandlerV1) FindMessageRoom(c *gin.Context) {
	appG := app.Gin{C: c}

	// authData, err := middleware.GetAuthFromCtx(c)
	// fmt.Println("auth ", authData.Roles)

	// params, err := utils.GetParamsFromRequest(c, model.Message{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *model.MessageRoomFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(params)
	Nodes, err := h.Services.MessageRoom.FindMessageRoom(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodes)
}

func (h *HandlerV1) UpdateMessageRoom(c *gin.Context) {
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
	var input *model.MessageRoom
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.MessageRoom.UpdateMessageRoom(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteMessageRoom(c *gin.Context) {
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

	node, err := h.Services.MessageRoom.DeleteMessageRoom(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateOrExistMessageRoom(c *gin.Context, input *model.MessageRoom) (*model.MessageRoom, error) {
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

	var result *model.MessageRoom

	// check exist product.
	// domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"_id", input.ProductID},
	// 	},
	// 	Options: domain.Options{
	// 		Limit: 1,
	// 	},
	// }
	prodID := input.ProductID.Hex()
	existNode, err := h.Services.Product.FindProduct(&model.ProductFilter{ID: []*string{&prodID}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	if len(existNode.Data) == 0 {
		//appG.ResponseError(http.StatusBadRequest, errors.New("not found node"), nil)
		return result, nil
	}

	input.TakeUserID = existNode.Data[0].UserID

	// // check exist message
	// existMessage, err := h.services.Message.FindMessage(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"product_id", input.ProductID},
	// 		{"user_id", userIDPrimitive},
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
	result, err = h.Services.MessageRoom.CreateMessageRoom(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	return result, nil
}
