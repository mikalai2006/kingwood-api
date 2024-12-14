package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerOffer(router *gin.RouterGroup) {
	offer := router.Group("/offer")
	offer.POST("/find", h.FindOffer)
	offer.POST("", h.CreateOffer)
	offer.POST("/list", h.CreateListOffer)
	offer.PATCH("/:id", h.UpdateOffer)
	offer.DELETE("/:id", h.DeleteOffer)
}

func (h *HandlerV1) CreateOffer(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	// var input *model.Offer
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	input, er := utils.BindJSON2[model.OfferInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Offer, err := h.CreateOrExistOffer(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Offer)
}

func (h *HandlerV1) CreateListOffer(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input []*model.OfferInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	result := []*model.Offer{}
	for i := range input {

		Offer, err := h.CreateOrExistOffer(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Offer)
		// nodeIDPrimitive, err := primitive.ObjectIDFromHex(input[i].NodeID)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// tagIDPrimitive, err := primitive.ObjectIDFromHex(input[i].TagID)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// // existOffer, err := h.services.Offer.FindOffer(domain.RequestParams{
		// // 	Options: domain.Options{Limit: 1},
		// // 	Filter:  bson.D{{"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}},
		// // })
		// // if err != nil {
		// // 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// // 	return
		// // }

		// existOffer, err := h.services.Offer.FindOffer(domain.RequestParams{
		// 	Options: domain.Options{Limit: 1},
		// 	Filter:  bson.D{{"data.value", input[i].Data.Value}, {"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}}, // {"tag_id", input.TagID},
		// })
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// if len(existOffer.Data) == 0 {
		// 	Offer, err := h.services.Offer.CreateOffer(userID, input[i])
		// 	if err != nil {
		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
		// 		return
		// 	}
		// 	result = append(result, Offer)
		// } else {
		// 	result = append(result, &existOffer.Data[0])
		// }
		// // else {
		// // 	fmt.Println("Exist data for ", existOffer.Data[0])
		// // }

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find Offers by params
// @Security ApiKeyAuth
// @Offers Offer
// @Description Input params for search Offers
// @ModuleID Offer
// @Accept  json
// @Produce  json
// @Param input query OfferInput true "params for search Offer"
// @Success 200 {object} []model.Offer
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Offer [get].
func (h *HandlerV1) FindOffer(c *gin.Context) {
	appG := app.Gin{C: c}

	// params, err := utils.GetParamsFromRequest(c, model.Offer{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// implementation roles for user.
	// roles, err := middleware.GetRoles(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }

	var input *model.OfferFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	items, err := h.Services.Offer.FindOffer(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *HandlerV1) DeleteOffer(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	Offer, err := h.Services.Offer.DeleteOffer(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// // remove Offer votes.
	// OfferVotes, err := h.services.OfferVote.FindOfferVote(domain.RequestParams{
	// 	Filter: bson.D{{"Offer_id", Offer.ID}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// if len(OfferVotes.Data) > 0 {
	// 	for i, _ := range OfferVotes.Data {
	// 		_, err := h.services.OfferVote.DeleteOfferVote(OfferVotes.Data[i].ID.Hex())
	// 		if err != nil {
	// 			appG.ResponseError(http.StatusBadRequest, err, nil)
	// 			return
	// 		}
	// 	}
	// }

	c.JSON(http.StatusOK, Offer)
}

func (h *HandlerV1) CreateOrExistOffer(c *gin.Context, input *model.OfferInput) (*model.Offer, error) {
	appG := app.Gin{C: c}

	var result *model.Offer

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return result, err
	}
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	productIDPrimitive, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	// check exist node
	// domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"_id", productIDPrimitive}},
	// }
	existProducts, err := h.Services.Product.FindProduct(&model.ProductFilter{ID: []*string{&input.ProductID}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	if len(existProducts.Data) == 0 {
		// appG.ResponseError(http.StatusBadRequest, model.ErrNodeNotFound, nil)
		return result, nil
	}
	input.UserProductID = existProducts.Data[0].UserID.Hex()

	// check exist Offer
	// domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"product_id", productIDPrimitive}, {"user_id", userIDPrimitive}}, //  {"tag_id", input.TagID},
	// }
	existOffer, err := h.Services.Offer.FindOffer(&model.OfferFilter{
		ProductID: []*primitive.ObjectID{&productIDPrimitive},
		UserID:    &userIDPrimitive,
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	// fmt.Println(existOffer)

	if len(existOffer.Data) > 0 {

		// cost := input.Cost - existOffer.Data[0].Cost
		// _, err = h.Services.User.SetBal(userID, -int(cost))
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return result, err
		}

		data := model.Offer{Status: input.Status, Cost: input.Cost}

		result, err = h.Services.Offer.UpdateOffer(existOffer.Data[0].ID.Hex(), userID, &data)

		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return result, err
		}
		return &existOffer.Data[0], nil
	} else {
		result, err = h.Services.Offer.CreateOffer(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return result, err
		}
	}

	return result, nil
}

func (h *HandlerV1) UpdateOffer(c *gin.Context) {

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
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[model.Node](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// fmt.Println(data)
	var input *model.Offer
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Offer.UpdateOffer(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	status := 1
	if *document.Take == status { //  && document.Give == &status
		_, err = h.Services.Product.UpdateProduct(document.ProductID.Hex(), userID, &model.Product{Status: -1})
		if err != nil {
			appG.ResponseError(http.StatusInternalServerError, err, nil)
			return
		}
	}

	c.JSON(http.StatusOK, document)
}
