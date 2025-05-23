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
)

func (h *HandlerV1) registerOrder(router *gin.RouterGroup) {
	route := router.Group("/order")
	route.POST("/find", h.FindOrder)
	route.POST("", h.CreateOrder)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateOrder)
	route.POST("/list", h.CreateOrderList)
	route.DELETE("/:id", h.DeleteOrder)
}

func (h *HandlerV1) CreateOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.Order
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Order, err := h.CreateOrExistOrder(c, input) //h.services.Order.CreateOrder(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Order)
}

func (h *HandlerV1) CreateOrderList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.Order
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Order
	for i := range input {
		Order, err := h.CreateOrExistOrder(c, input[i]) //h.services.Order.CreateOrder(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Order)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Order Get all Orders
// @Security ApiKeyAuth
// @Tags Order
// @Description get all Orders
// @ModuleID Order
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Order
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Order [get].
// func (h *HandlerV1) GetAllOrder(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	params, err := utils.GetParamsFromRequest(c, domain.Order{}, &h.i18n)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	Orders, err := h.Services.Order.GetAllOrder(params)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, Orders)
// }

// @Summary Find Orders by params
// @Security ApiKeyAuth
// @Tags Order
// @Description Input params for search Orders
// @ModuleID Order
// @Accept  json
// @Produce  json
// @Param input query OrderInput true "params for search Order"
// @Success 200 {object} []domain.Order
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Order [get].
func (h *HandlerV1) FindOrder(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.OrderInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.OrderFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Orders, err := h.Services.Order.FindOrder(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Orders)
}

func (h *HandlerV1) GetOrderByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateOrder(c *gin.Context) {
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
	data, er := utils.BindJSON2[domain.OrderInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	fmt.Println("order data: ", data.Name)

	document, err := h.Services.Order.UpdateOrder(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteOrder(c *gin.Context) {
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

	user, err := h.Services.Order.DeleteOrder(id, userID) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) CreateOrExistOrder(c *gin.Context, input *domain.Order) (*domain.Order, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.Order

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// если в данных есть number, проверяем на существование такого номера
	if input.Number != 0 {
		lim := 1
		existOrders, err := h.Services.Order.FindOrder(&domain.OrderFilter{Year: input.Year, Number: &input.Number, Limit: &lim})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return result, err
		}
		if len(existOrders.Data) > 0 {
			// fmt.Println("existOrders =")
			// appG.ResponseError(http.StatusBadRequest, domain.ErrExistNumberOrder, nil)
			return result, domain.ErrExistNumberOrder
			// return &existOrders.Data[0], nil
		}
	}

	result, err = h.Services.Order.CreateOrder(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}
