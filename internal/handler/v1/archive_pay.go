package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerArchivePay(router *gin.RouterGroup) {
	route := router.Group("/archive_pay")
	// route.POST("", h.CreateArchivePay)
	route.POST("/populate", h.FindArchivePay)
	route.DELETE("/:id", h.DeleteArchivePay)
}

// func (h *HandlerV1) CreateArchivePay(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	// userID, err := middleware.GetUID(c)
// 	// if err != nil {
// 	// 	// c.AbortWithError(http.StatusUnauthorized, err)
// 	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 	// 	return
// 	// }

// 	var input *domain.ArchivePay
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	ArchivePay, err := h.CreateOrExistArchivePay(c, input) //h.services.ArchivePay.CreateArchivePay(userID, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, ArchivePay)
// }

// func (h *HandlerV1) CreateArchivePayList(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil || userID == "" {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return
// 	}

// 	var input []*domain.ArchivePay
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	if len(input) == 0 {
// 		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
// 		return
// 	}

// 	var result []*domain.ArchivePay
// 	for i := range input {
// 		ArchivePay, err := h.CreateOrExistArchivePay(c, input[i]) //h.services.ArchivePay.CreateArchivePay(userID, input)
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return
// 		}
// 		result = append(result, ArchivePay)
// 	}

// 	c.JSON(http.StatusOK, result)
// }

// @Summary Find ArchivePays by params
// @Security ApiKeyAuth
// @Tags ArchivePay
// @Description Input params for search ArchivePays
// @ModuleID ArchivePay
// @Accept  json
// @Produce  json
// @Param input query ArchivePayInput true "params for search ArchivePay"
// @Success 200 {object} []domain.ArchivePay
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/ArchivePay [get].
func (h *HandlerV1) FindArchivePay(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.ArchivePayInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.ArchivePayFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	ArchivePays, err := h.Services.ArchivePay.FindArchivePay(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, ArchivePays)
}

func (h *HandlerV1) DeleteArchivePay(c *gin.Context) {
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

	user, err := h.Services.ArchivePay.DeleteArchivePay(id, userID)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// func (h *HandlerV1) CreateOrExistArchivePay(c *gin.Context, input *domain.ArchivePay) (*domain.ArchivePay, error) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return nil, err
// 	}
// 	var result *domain.ArchivePay

// 	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return result, err
// 	// }

// 	// existArchivePays, err := h.services.ArchivePay.FindArchivePay(domain.RequestParams{
// 	// 	Options: domain.Options{Limit: 1},
// 	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"userId", userIDPrimitive}},
// 	// })
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return result, err
// 	// }
// 	// if len(existArchivePays.Data) > 0 {
// 	// 	fmt.Println("existArchivePays =")
// 	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
// 	// 	return &existArchivePays.Data[0], nil
// 	// }

// 	result, err = h.Services.ArchivePay.CreateArchivePay(userID, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}
// 	return result, nil
// }
