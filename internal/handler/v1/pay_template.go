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
	"go.mongodb.org/mongo-driver/bson"
)

func (h *HandlerV1) registerPayTemplate(router *gin.RouterGroup) {
	route := router.Group("/pay_template")
	route.GET("", h.FindPayTemplate)
	route.POST("", h.SetUserFromRequest, h.CreatePayTemplate)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdatePayTemplate)
	route.DELETE("/:id", h.SetUserFromRequest, h.DeletePayTemplate)
}

func (h *HandlerV1) CreatePayTemplate(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[domain.PayTemplate](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// find exist PayTemplate.
	existPayTemplate, err := h.Services.PayTemplate.FindPayTemplate(domain.RequestParams{Filter: bson.D{
		{"name", data.Name},
	}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(existPayTemplate.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("exist PayTemplate"), nil)
		return
	}

	PayTemplate, err := h.Services.PayTemplate.CreatePayTemplate(userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, PayTemplate)
}

// @Summary Find PayTemplates by params
// @Security ApiKeyAuth
// @PayTemplates PayTemplate
// @Description Input params for search PayTemplates
// @ModuleID PayTemplate
// @Accept  json
// @Produce  json
// @Param input query PayTemplateInput true "params for search PayTemplate"
// @Success 200 {object} []domain.PayTemplate
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/PayTemplate [get].
func (h *HandlerV1) FindPayTemplate(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.PayTemplateInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	PayTemplates, err := h.Services.PayTemplate.FindPayTemplate(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, PayTemplates)
}

func (h *HandlerV1) UpdatePayTemplate(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// // var input domain.PayTemplateInput
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
	// data, er := utils.BindJSON[domain.PayTemplate](a)
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
	data, er := utils.BindJSON2[domain.PayTemplateInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.PayTemplate.UpdatePayTemplate(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeletePayTemplate(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.Services.PayTemplate.DeletePayTemplate(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
