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

func (h *HandlerV1) registerTaskMontajWorker(router *gin.RouterGroup) {
	route := router.Group("/task_montaj_worker")
	route.POST("", h.CreateTaskMontajWorker)
	route.GET("", h.FindTaskMontajWorker)
	route.POST("/populate", h.FindTaskMontajWorkerPopulate)
	route.PATCH("/:id", h.SetUserFromRequest, h.UpdateTaskMontajWorker)
	route.POST("/list", h.CreateTaskMontajWorkerList)
	route.DELETE("/:id", h.DeleteTaskMontajWorker)
}

func (h *HandlerV1) CreateTaskMontajWorker(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *domain.TaskMontajWorker
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	TaskMontajWorker, err := h.CreateOrExistTaskMontajWorker(c, input) //h.services.TaskMontajWorker.CreateTaskMontajWorker(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, TaskMontajWorker)
}

func (h *HandlerV1) CreateTaskMontajWorkerList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.TaskMontajWorker
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.TaskMontajWorker
	for i := range input {
		TaskMontajWorker, err := h.CreateOrExistTaskMontajWorker(c, input[i]) //h.services.TaskMontajWorker.CreateTaskMontajWorker(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, TaskMontajWorker)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find TaskMontajWorkers by params
// @Security ApiKeyAuth
// @Tags TaskMontajWorker
// @Description Input params for search TaskMontajWorkers
// @ModuleID TaskMontajWorker
// @Accept  json
// @Produce  json
// @Param input query TaskMontajWorkerInput true "params for search TaskMontajWorker"
// @Success 200 {object} []domain.TaskMontajWorker
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/TaskMontajWorker [get].
func (h *HandlerV1) FindTaskMontajWorker(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.TaskMontajWorkerInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	TaskMontajWorkers, err := h.Services.TaskMontajWorker.FindTaskMontajWorker(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, TaskMontajWorkers)
}

func (h *HandlerV1) FindTaskMontajWorkerPopulate(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.TaskMontajWorkerInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.TaskMontajWorkerFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	TaskMontajWorkers, err := h.Services.TaskMontajWorker.FindTaskMontajWorkerPopulate(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, TaskMontajWorkers)
}

func (h *HandlerV1) UpdateTaskMontajWorker(c *gin.Context) {
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
	data, er := utils.BindJSON2[domain.TaskMontajWorkerInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.TaskMontajWorker.UpdateTaskMontajWorker(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) CreateOrExistTaskMontajWorker(c *gin.Context, input *domain.TaskMontajWorker) (*domain.TaskMontajWorker, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *domain.TaskMontajWorker

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// existTaskMontajWorkers, err := h.services.TaskMontajWorker.FindTaskMontajWorker(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"user_id", userIDPrimitive}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }
	// if len(existTaskMontajWorkers.Data) > 0 {
	// 	fmt.Println("existTaskMontajWorkers =")
	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
	// 	return &existTaskMontajWorkers.Data[0], nil
	// }

	result, err = h.Services.TaskMontajWorker.CreateTaskMontajWorker(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}

func (h *HandlerV1) DeleteTaskMontajWorker(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.Services.TaskMontajWorker.DeleteTaskMontajWorker(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
