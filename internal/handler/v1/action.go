package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerAction(router *gin.RouterGroup) {
	action := router.Group("/action")
	action.GET("/", h.FindAction)
	action.POST("/", h.CreateAction)
	action.POST("/list/", h.CreateListAction)
	action.PATCH("/:id", h.UpdateAction)
	action.DELETE("/:id", h.DeleteAction)
}

func (h *HandlerV1) CreateAction(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	// var input *model.Action
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage // map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	input, er := utils.BindJSON2[model.ActionInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	serviceIDPrimitive, err := primitive.ObjectIDFromHex(input.ServiceID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	existAction, err := h.Services.Action.FindAction(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"type", input.Type}, {"service_id", serviceIDPrimitive}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if len(existAction.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, model.ErrActionExistValue, nil)
		return
	}

	Action, err := h.Services.Action.CreateAction(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Action)
}

func (h *HandlerV1) CreateListAction(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.ActionInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	result := []*model.Action{}
	for i := range input {
		serviceIDPrimitive, err := primitive.ObjectIDFromHex(input[i].ServiceID)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		existAction, err := h.Services.Action.FindAction(domain.RequestParams{
			Options: domain.Options{Limit: 1},
			Filter:  bson.D{{"type", input[i].Type}, {"service_id", serviceIDPrimitive}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		if len(existAction.Data) == 0 {
			Action, err := h.Services.Action.CreateAction(userID, input[i])
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			result = append(result, Action)
		}

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Action Get all Actions
// @Security ApiKeyAuth
// @Actions Action
// @Description get all Actions
// @ModuleID Action
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Action
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Action [get].
func (h *HandlerV1) GetAllAction(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Action{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Actions, err := h.Services.Action.GetAllAction(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Actions)
}

// @Summary Find Actions by params
// @Security ApiKeyAuth
// @Actions Action
// @Description Input params for search Actions
// @ModuleID Action
// @Accept  json
// @Produce  json
// @Param input query ActionInput true "params for search Action"
// @Success 200 {object} []model.Action
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Action [get].
func (h *HandlerV1) FindAction(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.ActionInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Actions, err := h.Services.Action.FindAction(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Actions)
}

func (h *HandlerV1) GetActionByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateAction(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	var a map[string]json.RawMessage // map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[model.ActionInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.Services.Action.UpdateAction(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteAction(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.Services.Action.DeleteAction(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
