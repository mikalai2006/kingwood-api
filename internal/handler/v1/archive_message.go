package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

func (h *HandlerV1) registerArchiveMessage(router *gin.RouterGroup) {
	ArchiveMessage := router.Group("/archive_message")
	ArchiveMessage.POST("/find", h.FindArchiveMessage)
	ArchiveMessage.DELETE("/:id", h.DeleteArchiveMessage)
}

func (h *HandlerV1) FindArchiveMessage(c *gin.Context) {
	appG := app.Gin{C: c}

	// authData, err := middleware.GetAuthFromCtx(c)
	// fmt.Println("auth ", authData.Roles)

	// params, err := utils.GetParamsFromRequest(c, domain.ArchiveMessage{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input *domain.ArchiveMessageFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(params)
	result, err := h.Services.ArchiveMessage.FindArchiveMessage(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *HandlerV1) DeleteArchiveMessage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	// implementation roles for user.
	// roles, err := middleware.GetRoles(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// if !utils.Contains(roles, "admin") {
	// 	appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
	// 	return
	// }

	node, err := h.Services.ArchiveMessage.DeleteArchiveMessage(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}
