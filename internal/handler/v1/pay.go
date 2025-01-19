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

func (h *HandlerV1) registerPay(router *gin.RouterGroup) {
	route := router.Group("/pay")
	route.POST("", h.CreatePay)
	route.POST("/populate", h.FindPay)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdatePay)
	// route.POST("/list", h.CreatePayList)
	route.DELETE("/:id", h.DeletePay)
}

func (h *HandlerV1) CreatePay(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.Pay
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Pay, err := h.CreateOrExistPay(c, input) //h.services.Pay.CreatePay(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Pay)
}

func (h *HandlerV1) CreatePayList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.Pay
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Pay
	for i := range input {
		Pay, err := h.CreateOrExistPay(c, input[i]) //h.services.Pay.CreatePay(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Pay)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find Pays by params
// @Security ApiKeyAuth
// @Tags Pay
// @Description Input params for search Pays
// @ModuleID Pay
// @Accept  json
// @Produce  json
// @Param input query PayInput true "params for search Pay"
// @Success 200 {object} []domain.Pay
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Pay [get].
func (h *HandlerV1) FindPay(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.PayInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.PayFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Pays, err := h.Services.Pay.FindPay(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Pays)
}

func (h *HandlerV1) GetPayByID(c *gin.Context) {

}

func (h *HandlerV1) UpdatePay(c *gin.Context) {
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
	data, er := utils.BindJSON2[domain.PayInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Pay.UpdatePay(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeletePay(c *gin.Context) {
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

	user, err := h.Services.Pay.DeletePay(id, userID)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) CreateOrExistPay(c *gin.Context, input *domain.Pay) (*domain.Pay, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.Pay

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// existPays, err := h.services.Pay.FindPay(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"user_id", userIDPrimitive}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }
	// if len(existPays.Data) > 0 {
	// 	fmt.Println("existPays =")
	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
	// 	return &existPays.Data[0], nil
	// }

	result, err = h.Services.Pay.CreatePay(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}
