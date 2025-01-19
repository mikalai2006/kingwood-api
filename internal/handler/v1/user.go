package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) RegisterUser(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.POST("", h.CreateUser)
	user.POST("/populate", h.FindUser)
	user.GET("/:id", h.GetUser)
	user.DELETE("/:id", h.DeleteUser)
	user.PATCH("/:id", h.UpdateUser)
}

// @Summary Get user by Id
// @Tags user
// @Description get user info
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [get].
func (h *HandlerV1) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	user, err := h.Services.User.GetUser(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// // get auth data for user
	// authData, err := h.services.GetAuth(user.UserID.Hex())
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// if !authData.ID.IsZero() {
	// 	user.Md = authData.MaxDistance
	// 	user.Roles = authData.Roles
	// }

	c.JSON(http.StatusOK, user)
}

// type InputUser struct {
// 	domain.RequestParams
// 	domain.User
// }

// @Summary Find few users
// @Security ApiKeyAuth
// @Tags user
// @Description Input params for search users
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param input query domain.UserInput true "params for search users"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user [get].
func (h *HandlerV1) FindUser(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, domain.UserInput{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.UserFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	users, err := h.Services.User.FindUser(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *HandlerV1) CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.User
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, er := h.Services.User.CreateUser(userID, input)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Delete user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [delete].
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	// var input domain.User
	// if err := c.BindJSON(&input); err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())

	// 	return
	// }

	user, err := h.Services.User.DeleteUser(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Security ApiKeyAuth
// @Tags user
// @Description Update user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body domain.User true "body for update user"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [put].
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input domain.UserInput
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, err := h.Services.User.GetUser(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var imageInput = &domain.ImageInput{}
	imageInput.Service = "user"
	imageInput.ServiceID = user.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "user"

	paths, err := utils.UploadResizeMultipleFile(c, imageInput, "images", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	var images []domain.Image
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		image, err := h.Services.Image.CreateImage(userID, imageInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		images = append(images, image)
	}

	result, err := h.Services.User.UpdateUser(id, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// limit := 1
	// result, err := h.Services.User.FindUser(&domain.UserFilter{ID: []string{user.ID.Hex()}, Limit: &limit})
	// // domain.RequestParams{Filter: bson.D{{"_id", user.ID}}})
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	c.JSON(http.StatusOK, result)
}
