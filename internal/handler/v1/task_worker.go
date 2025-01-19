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

func (h *HandlerV1) registerTaskWorker(router *gin.RouterGroup) {
	route := router.Group("/task_worker")
	route.POST("", h.CreateTaskWorker)
	// route.GET("", h.FindTaskWorker)
	route.POST("/populate", h.FindTaskWorkerPopulate)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateTaskWorker)
	route.POST("/list", h.CreateTaskWorkerList)
	route.DELETE("/:id", h.DeleteTaskWorker)
}

func (h *HandlerV1) CreateTaskWorker(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.TaskWorker
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	TaskWorker, err := h.CreateOrExistTaskWorker(c, input) //h.services.TaskWorker.CreateTaskWorker(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, TaskWorker)
}

func (h *HandlerV1) CreateTaskWorkerList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.TaskWorker
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.TaskWorker
	for i := range input {
		TaskWorker, err := h.CreateOrExistTaskWorker(c, input[i]) //h.services.TaskWorker.CreateTaskWorker(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, TaskWorker)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find TaskWorkers by params
// @Security ApiKeyAuth
// @Tags TaskWorker
// @Description Input params for search TaskWorkers
// @ModuleID TaskWorker
// @Accept  json
// @Produce  json
// @Param input query TaskWorkerInput true "params for search TaskWorker"
// @Success 200 {object} []domain.TaskWorker
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/TaskWorker [get].
// func (h *HandlerV1) FindTaskWorker(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	params, err := utils.GetParamsFromRequest(c, domain.TaskWorkerInputData{}, &h.i18n)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	TaskWorkers, err := h.Services.TaskWorker.FindTaskWorker(params)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

//		c.JSON(http.StatusOK, TaskWorkers)
//	}
func (h *HandlerV1) FindTaskWorkerPopulate(c *gin.Context) {
	appG := app.Gin{C: c}
	// params, err := utils.GetParamsFromRequest(c, domain.TaskWorkerInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.TaskWorkerFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	TaskWorkers, err := h.Services.TaskWorker.FindTaskWorkerPopulate(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, TaskWorkers)
}

func (h *HandlerV1) UpdateTaskWorker(c *gin.Context) {
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
	data, er := utils.BindJSON2[domain.TaskWorkerInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.TaskWorker.UpdateTaskWorker(id, userID, &data, 1)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) CreateOrExistTaskWorker(c *gin.Context, input *domain.TaskWorker) (*domain.TaskWorker, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.TaskWorker

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// existTaskWorkers, err := h.services.TaskWorker.FindTaskWorker(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"user_id", userIDPrimitive}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }
	// if len(existTaskWorkers.Data) > 0 {
	// 	fmt.Println("existTaskWorkers =")
	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
	// 	return &existTaskWorkers.Data[0], nil
	// }

	result, err = h.Services.TaskWorker.CreateTaskWorker(userID, input, 1)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}

func (h *HandlerV1) DeleteTaskWorker(c *gin.Context) {
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

	user, err := h.Services.TaskWorker.DeleteTaskWorker(id, userID) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
