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

func (h *HandlerV1) registerWorkHistory(router *gin.RouterGroup) {
	route := router.Group("/work_history")
	route.POST("", h.CreateWorkHistory)
	route.POST("/find", h.FindWorkHistory)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateWorkHistory)
	route.POST("/list", h.CreateWorkHistoryList)
}

func (h *HandlerV1) CreateWorkHistory(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.WorkHistory
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Work, err := h.CreateOrExistWorkHistory(c, input) //h.services.Work.CreateWork(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Work)
}

func (h *HandlerV1) CreateWorkHistoryList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.WorkHistory
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.WorkHistory
	for i := range input {
		Work, err := h.CreateOrExistWorkHistory(c, input[i]) //h.services.Work.CreateWork(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Work)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find Works by params
// @Security ApiKeyAuth
// @Tags Work
// @Description Input params for search Works
// @ModuleID Work
// @Accept  json
// @Produce  json
// @Param input query WorkInput true "params for search Work"
// @Success 200 {object} []domain.Work
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Work [get].
func (h *HandlerV1) FindWorkHistory(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.WorkHistoryInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.WorkHistoryFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Works, err := h.Services.WorkHistory.FindWorkHistory(*input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Works)
}

func (h *HandlerV1) UpdateWorkHistory(c *gin.Context) {
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
	data, er := utils.BindJSON2[domain.WorkHistoryInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.WorkHistory.UpdateWorkHistory(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteWorkHistory(c *gin.Context) {

}

func (h *HandlerV1) CreateOrExistWorkHistory(c *gin.Context, input *domain.WorkHistory) (*domain.WorkHistory, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.WorkHistory

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// existWorks, err := h.services.Work.FindWork(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"user_id", userIDPrimitive}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }
	// if len(existWorks.Data) > 0 {
	// 	fmt.Println("existWorks =")
	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
	// 	return &existWorks.Data[0], nil
	// }

	result, err = h.Services.WorkHistory.CreateWorkHistory(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}
