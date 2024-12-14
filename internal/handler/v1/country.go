package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *HandlerV1) RegisterCountry(router *gin.RouterGroup) {
	country := router.Group("/country")
	country.POST("", h.SetUserFromRequest, h.createCountry)
	country.GET("", h.findCountry)
	country.GET("/:id", h.getCountry)
	country.POST("/list/", h.SetUserFromRequest, h.CreateListCountry)
	country.PATCH("/:id", h.SetUserFromRequest, h.updateCountry)
	country.DELETE("/:id", h.SetUserFromRequest, h.deleteCountry)
}

func (h *HandlerV1) createCountry(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.CountryInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Country.CreateCountry(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) CreateListCountry(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.CountryInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []domain.Country
	for i := range input {
		existCountry, err := h.Services.Country.FindCountry(domain.RequestParams{
			Options: domain.Options{Limit: 1},
			Filter:  bson.D{{"code", input[i].Code}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		// existLatLon := false
		// if len(existCountry.Data) > 0 {
		// 	existLatLon = input[i].Lat == existCountry.Data[0].Lat && input[i].Lon == existCountry.Data[0].Lon
		// 	progress := 0
		// 	if existLatLon {
		// 		progress = 100
		// 	}

		// 	_, err := h.services.Ticket.CreateTicket(userID, &model.Ticket{
		// 		Title:       "Double osm object",
		// 		Description: fmt.Sprintf("[osmId]%s[/osmId]: [coords]%v,%v[/coords], [existCoords]%v,%v[/existCoords]", input[i].OsmID, input[i].Lat, input[i].Lon, existCountry.Data[0].Lat, existCountry.Data[0].Lon),
		// 		Status:      !existLatLon,
		// 		Progress:    progress,
		// 	})
		// 	if err != nil {
		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
		// 		return
		// 	}
		// 	// fmt.Println("Double node:::", input[i].OsmID, input[i].Lat, input[i].Lon)
		// }
		if len(existCountry.Data) == 0 {
			country, err := h.Services.Country.CreateCountry(userID, input[i])
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			result = append(result, country)
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *HandlerV1) getCountry(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.Services.Country.GetCountry(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) findCountry(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.CountryInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	documents, err := h.Services.Country.FindCountry(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *HandlerV1) updateCountry(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.CountryInput
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.Services.Country.UpdateCountry(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) deleteCountry(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}
	// var input domain.Page
	// if err := c.BindJSON(&input); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)

	// 	return
	// }

	document, err := h.Services.Country.DeleteCountry(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
