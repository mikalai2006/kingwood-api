package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerQuestion(router *gin.RouterGroup) {
	question := router.Group("/question")
	question.POST("/", h.SetUserFromRequest, h.CreateQuestion)
	question.POST("/find", h.FindQuestion)
	// question.POST("/list/", h.SetUserFromRequest, h.CreateListQuestion)
	question.PATCH("/:id", h.SetUserFromRequest, h.UpdateQuestion)
	question.DELETE("/:id", h.SetUserFromRequest, h.DeleteQuestion)
}

func (h *HandlerV1) CreateQuestion(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	// var input *model.Question
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage // map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	input, er := utils.BindJSON2[model.QuestionInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	productID := input.ProductID.Hex()
	existProducts, err := h.Services.Product.FindProduct(&model.ProductFilter{ID: []*string{&productID}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(existProducts.Data) == 0 {
		// appG.ResponseError(http.StatusBadRequest, model.ErrNodeNotFound, nil)
		return
	}
	input.UserProductID = existProducts.Data[0].UserID

	// existQuestion, err := h.services.Question.FindQuestion(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"value", input.ProductID}}, // {"tag_id", input.TagID},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// if len(existQuestion.Data) > 0 {
	// 	appG.ResponseError(http.StatusBadRequest, model.ErrQuestionExistValue, nil)
	// 	return
	// }

	Question, err := h.Services.Question.CreateQuestion(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Question)
}

// func (h *HandlerV1) CreateListQuestion(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return
// 	}

// 	var input []*model.QuestionInput
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	if len(input) == 0 {
// 		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
// 		return
// 	}

// 	result := []*model.Question{}
// 	for i := range input {
// 		// existQuestion, err := h.services.Question.FindQuestion(domain.RequestParams{
// 		// 	Options: domain.Options{Limit: 1},
// 		// 	Filter:  bson.D{{"value", input[i].Value}, {"tag_id", input[i].TagID}},
// 		// })
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return
// 		// }

// 		// if len(existQuestion.Data) == 0 {
// 		Question, err := h.Services.Question.CreateQuestion(userID, input[i])
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return
// 		}
// 		result = append(result, Question)
// 		// }

// 	}

// 	c.JSON(http.StatusOK, result)
// }

// @Summary Find Questions by params
// @Security ApiKeyAuth
// @Questions Question
// @Description Input params for search Questions
// @ModuleID Question
// @Accept  json
// @Produce  json
// @Param input query QuestionInput true "params for search Question"
// @Success 200 {object} []model.Question
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Question [get].
func (h *HandlerV1) FindQuestion(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, model.QuestionInput{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *model.QuestionFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	fmt.Println("input: ", input)
	Questions, err := h.Services.Question.FindQuestion(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Questions)
}

func (h *HandlerV1) GetQuestionByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateQuestion(c *gin.Context) {
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
	var a map[string]json.RawMessage // map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[model.QuestionInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.Services.Question.UpdateQuestion(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteQuestion(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.Services.Question.DeleteQuestion(id, userID) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
