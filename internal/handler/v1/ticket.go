package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerTicket(router *gin.RouterGroup) {
	ticket := router.Group("/ticket")
	ticket.GET("", h.FindTicket)
	ticket.POST("", h.CreateTicket)
	ticket.POST("/:id/message", h.CreateTicketMessage)
	ticket.POST("/list/", h.CreateListTicket)
	ticket.DELETE("/:id", h.DeleteTicket)
}

func (h *HandlerV1) CreateTicket(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *model.Ticket
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println("input: ", input)

	Ticket, err := h.Services.Ticket.CreateTicket(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Ticket)
}

func (h *HandlerV1) CreateTicketMessage(c *gin.Context) {
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
		appG.ResponseError(http.StatusBadRequest, errors.New("need id"), nil)
		return
	}

	var input model.TicketMessage
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	ticketIdPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	// fmt.Println("input: ", input)

	tickets, err := h.Services.Ticket.FindTicket(domain.RequestParams{Filter: bson.D{
		{"_id", ticketIdPrimitive},
	}, Options: domain.Options{Limit: 1}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(tickets.Data) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("not found ticket"), nil)
		return
	}

	// fmt.Println("ticket: ", tickets.Data[0])
	ticketMessage, err := h.Services.Ticket.CreateTicketMessage(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var imageInput = &domain.ImageInput{}
	imageInput.Service = "ticket"
	imageInput.ServiceID = ticketMessage.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "ticket"

	paths, err := utils.UploadResizeMultipleFile(c, imageInput, "images", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}

	var result []domain.Image
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		image, err := h.Services.Image.CreateImage(userID, imageInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, image)
	}

	// fmt.Println("paths: ", paths)

	c.JSON(http.StatusOK, ticketMessage)
}

func (h *HandlerV1) CreateListTicket(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.Ticket
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.Ticket
	for i := range input {
		Ticket, err := h.Services.Ticket.CreateTicket(userID, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Ticket)

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Ticket Get all Tickets
// @Security ApiKeyAuth
// @Tickets Ticket
// @Description get all Tickets
// @ModuleID Ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Ticket
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Ticket [get].
func (h *HandlerV1) GetAllTicket(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Ticket{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tickets, err := h.Services.Ticket.GetAllTicket(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tickets)
}

// @Summary Find Tickets by params
// @Security ApiKeyAuth
// @Tickets Ticket
// @Description Input params for search Tickets
// @ModuleID Ticket
// @Accept  json
// @Produce  json
// @Param input query TicketInput true "params for search Ticket"
// @Success 200 {object} []model.Ticket
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Ticket [get].
func (h *HandlerV1) FindTicket(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.TicketInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tickets, err := h.Services.Ticket.FindTicket(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tickets)
}

func (h *HandlerV1) GetTicketByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateTicket(c *gin.Context) {

}

func (h *HandlerV1) DeleteTicket(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.Services.Ticket.DeleteTicket(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
