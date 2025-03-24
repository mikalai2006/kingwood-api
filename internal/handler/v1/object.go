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

func (h *HandlerV1) registerObject(router *gin.RouterGroup) {
	route := router.Group("/object")
	route.POST("", h.CreateObject)
	route.POST("/find", h.FindObject)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateObject)
	route.POST("/list", h.CreateObjectList)
	route.DELETE("/:id", h.DeleteObject)
}

func (h *HandlerV1) CreateObject(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.Object
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Object, err := h.CreateOrExistObject(c, input) //h.services.Object.CreateObject(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Object)
}

func (h *HandlerV1) CreateObjectList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.Object
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Object
	for i := range input {
		Object, err := h.CreateOrExistObject(c, input[i]) //h.services.Object.CreateObject(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Object)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find Objects by params
// @Security ApiKeyAuth
// @Tags Object
// @Description Input params for search Objects
// @ModuleID Object
// @Accept  json
// @Produce  json
// @Param input query ObjectInput true "params for search Object"
// @Success 200 {object} []domain.Object
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Object [get].
func (h *HandlerV1) FindObject(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.ObjectInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.ObjectFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Objects, err := h.Services.Object.FindObject(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Objects)
}

func (h *HandlerV1) GetObjectByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateObject(c *gin.Context) {
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
	data, er := utils.BindJSON2[domain.ObjectInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Object.UpdateObject(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteObject(c *gin.Context) {
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
	}

	user, err := h.Services.Object.DeleteObject(id, userID)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) CreateOrExistObject(c *gin.Context, input *domain.Object) (*domain.Object, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.Object

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// existObjects, err := h.services.Object.FindObject(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"userId", userIDPrimitive}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }
	// if len(existObjects.Data) > 0 {
	// 	fmt.Println("existObjects =")
	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
	// 	return &existObjects.Data[0], nil
	// }

	result, err = h.Services.Object.CreateObject(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}
