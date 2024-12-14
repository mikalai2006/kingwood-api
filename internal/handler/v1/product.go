package v1

import (
	"errors"
	"fmt"
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

func (h *HandlerV1) registerProduct(router *gin.RouterGroup) {
	node := router.Group("/product")
	node.POST("/find", h.FindProduct)
	node.POST("/myoffers", h.FindProductForMyOffers)
	node.POST("", h.CreateProduct)
	node.PATCH("/:id", h.UpdateProduct)
	node.DELETE("/:id", h.DeleteProduct)
	// node.POST("/list/", h.CreateListNode)
}

func (h *HandlerV1) CreateProduct(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	// var input *model.Product
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }

	// // node, err := h.CreateOrExistProduct(c, input)
	// // if err != nil {
	// // 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// // 	return
	// // }
	// node, err := h.services.Product.CreateProduct(userID, input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var input model.ProductInputData
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	fmt.Println("CreateProduct input: ", input)

	node, err := h.Services.Product.CreateProduct(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var imageInput = &domain.ImageInput{}
	imageInput.Service = "product"
	imageInput.ServiceID = node.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "product"

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

	node.Images = result

	c.JSON(http.StatusOK, node)
}

// @Summary Find Nodes by params
// @Security ApiKeyAuth
// @Tags Node
// @Description Input params for search Nodes
// @ModuleID Node
// @Accept  json
// @Produce  json
// @Param input query NodeInput true "params for search Node"
// @Success 200 {object} []domain.Node
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Node [get].
func (h *HandlerV1) FindProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	// authData, err := middleware.GetAuthFromCtx(c)
	// fmt.Println("auth ", authData.Roles)
	var input *model.ProductFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// params, err := utils.GetParamsFromRequest(c, model.ProductFilter{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	items, err := h.Services.Product.FindProduct(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *HandlerV1) FindProductForMyOffers(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// domain.RequestParams{Filter: bson.D{{"owner_id", userID}}, Options: domain.Options{Limit: 100}}
	offers, err := h.Services.Offer.FindOffer(&model.OfferFilter{
		UserID: &userIDPrimitive,
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	productIDs := make([]*primitive.ObjectID, len(offers.Data))
	for i := range offers.Data {
		// productIDPrimitive, err := primitive.ObjectIDFromHex(*params.ProductID[i])
		// if err != nil {
		// 	return response, err
		// }
		productIDs[i] = &offers.Data[i].ProductID
	}

	items, err := h.Services.Product.FindProduct(&model.ProductFilter{ProductID: productIDs})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *HandlerV1) UpdateProduct(c *gin.Context) {
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
	var input *model.Product
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.Services.Product.UpdateProduct(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	// implementation roles for user.
	roles, err := middleware.GetRoles(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if !utils.Contains(roles, "admin") {
		appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
		return
	}

	node, err := h.Services.Product.DeleteProduct(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// find all images for remove.
	images, err := h.Services.Image.FindImage(domain.RequestParams{
		Filter: bson.D{
			{"service", "node"},
			{"service_id", node.ID.Hex()},
		},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	for i, _ := range images.Data {
		_, err := h.Services.Image.DeleteImage(images.Data[i].ID.Hex())
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// fmt.Println("Remove image", image.ID)
	}

	// // find all nodedata for remove.
	// nodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"node_id", node.ID},
	// 	},
	// 	Options: domain.Options{Limit: 10000},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// for i, _ := range nodedata.Data {
	// 	_, err := h.services.Nodedata.DeleteNodedata(nodedata.Data[i].ID.Hex())
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return
	// 	}
	// 	// fmt.Println("Remove nodedata: ", nodedata.Data[i].ID.Hex())
	// }

	// // find all reviews for remove.
	// reviews, err := h.services.Review.FindReview(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"node_id", node.ID},
	// 	},
	// 	Options: domain.Options{Limit: 10000},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// for i, _ := range reviews.Data {
	// 	_, err := h.services.Review.DeleteReview(reviews.Data[i].ID.Hex())
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return
	// 	}
	// 	// fmt.Println("Remove review: ", reviews.Data[i].ID.Hex())
	// }

	// // find all audits for remove.
	// nodeaudits, err := h.services.NodeAudit.FindNodeAudit(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"node_id", node.ID},
	// 	},
	// 	Options: domain.Options{Limit: 10000},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// for i, _ := range nodeaudits.Data {
	// 	_, err := h.services.NodeAudit.DeleteNodeAudit(nodeaudits.Data[i].ID.Hex())
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return
	// 	}
	// 	// fmt.Println("Remove nodeaudits: ", nodeaudits.Data[i].ID.Hex())
	// }

	// // find all vote of node for remove.
	// nodeVotes, err := h.services.NodeVote.FindNodeVote(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"node_id", node.ID},
	// 	},
	// 	Options: domain.Options{Limit: 10000},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// for i, _ := range nodeVotes.Data {
	// 	_, err := h.services.NodeVote.DeleteNodeVote(nodeVotes.Data[i].ID.Hex())
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return
	// 	}
	// 	// fmt.Println("Remove nodeaudits: ", nodeaudits.Data[i].ID.Hex())
	// }

	// // Remove address.
	// nodeAlsoOsmID, err := h.services.Node.FindNode(domain.RequestParams{
	// 	Filter: bson.D{{"osm_id", node.OsmID}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// if len(nodeAlsoOsmID.Data) == 0 {
	// 	addr, err := h.services.Address.FindAddress(domain.RequestParams{
	// 		Filter: bson.D{{"osm_id", node.OsmID}},
	// 	})
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return
	// 	}
	// 	if len(addr.Data) > 0 {
	// 		_, err = h.services.Address.DeleteAddress(addr.Data[0].ID.Hex())
	// 		if err != nil {
	// 			appG.ResponseError(http.StatusBadRequest, err, nil)
	// 			return
	// 		}
	// 	}
	// }

	c.JSON(http.StatusOK, node)
}

// func (h *HandlerV1) CreateOrExistProduct(c *gin.Context, input *model.Product) (*model.Product, error) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return nil, err
// 	}

// 	// // check exist node
// 	// existNode, err := h.services.Node.FindNode(domain.RequestParams{
// 	// 	Filter: bson.D{
// 	// 		{"lat", bson.M{"$lt": node.Lat + 0.001, "$gt": node.Lat - 0.001}},
// 	// 		{"lon", bson.M{"$lt": node.Lon + 0.001, "$gt": node.Lon - 0.001}},
// 	// 		{"type", node.Type},
// 	// 	},
// 	// 	Options: domain.Options{
// 	// 		Limit: 1,
// 	// 	},
// 	// })
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// if len(existNode.Data) > 0 {
// 	// 	return &existNode.Data[0], nil
// 	// }
// 	// // check exist node
// 	// existNode, err := h.services.Node.FindNode(domain.RequestParams{
// 	// 	Filter: bson.D{
// 	// 		{"lat", bson.M{"$lt": input.Lat + 0.00015, "$gt": input.Lat - 0.00015}},
// 	// 		{"lon", bson.M{"$lt": input.Lon + 0.00015, "$gt": input.Lon - 0.00015}},
// 	// 		{"type", input.Type},
// 	// 	},
// 	// 	Options: domain.Options{
// 	// 		Limit: 1,
// 	// 	},
// 	// })
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return nil, err
// 	// }

// 	// // if exist node
// 	// if len(existNode.Data) > 0 {
// 	// 	//appG.ResponseError(http.StatusBadRequest, errors.New("existSameNode"), nil)
// 	// 	return &existNode.Data[0], nil
// 	// }

// 	// // Get address.
// 	// pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", input.Lat, input.Lon))
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return nil, err
// 	// }
// 	// r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
// 	// r.Header.Add("User-Agent", "a127.0.0.1")

// 	// resp, err := http.DefaultClient.Do(r)
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return nil, err
// 	// }
// 	// defer resp.Body.Close()

// 	// bytes, err := io.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return nil, err
// 	// }
// 	// var bodyResponse domain.ResponseNominatim
// 	// if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, e, nil)
// 	// 	return nil, err
// 	// }

// 	var result *model.Product

// 	if bodyResponse.OsmID != 0 {
// 		// Check address in to db
// 		osmID := fmt.Sprintf("%s/%d", bodyResponse.OsmType, bodyResponse.OsmID)
// 		adrDB, err := h.services.Address.FindAddress(domain.RequestParams{Options: domain.Options{Limit: 1},
// 			Filter: bson.D{{"osm_id", osmID}}})
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return nil, err
// 		}

// 		address := &domain.Address{}
// 		if len(adrDB.Data) > 0 {
// 			address = &adrDB.Data[0]
// 		} else {
// 			address, err = h.services.Address.CreateAddress(userID, &domain.AddressInput{
// 				OsmID:    osmID,
// 				Address:  bodyResponse.Address,
// 				DAddress: bodyResponse.DisplayName,
// 			})
// 			if err != nil {
// 				appG.ResponseError(http.StatusBadRequest, err, nil)
// 				return nil, err
// 			}
// 		}

// 		input.OsmID = address.OsmID

// 		if bodyResponse.Name == "" {
// 			arrStr := strings.Split(address.DAddress, ",")
// 			nameNode := ""
// 			if len(arrStr) >= 2 {
// 				nameNode = fmt.Sprintf("%s, %s", arrStr[1], arrStr[0])
// 			} else {
// 				nameNode = arrStr[0]
// 			}
// 			input.Name = strings.TrimSpace(nameNode)
// 		} else {
// 			input.Name = bodyResponse.Name
// 		}

// 		if ccode, ok := bodyResponse.Address["country_code"]; ok {
// 			input.CCode = ccode.(string)
// 		}

// 		node, err := h.services.Node.CreateNode(userID, input)
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return nil, err
// 		}

// 		if len(input.Data) > 0 {
// 			for i := range input.Data {
// 				inputNodedata := &model.NodedataInput{
// 					NodeID:   node.ID.Hex(),
// 					Data:     input.Data[i].Data,
// 					TagID:    input.Data[i].TagID.Hex(),
// 					TagoptID: input.Data[i].TagoptID.Hex(),
// 				}

// 				Nodedata, err := h.CreateOrExistNodedata(c, inputNodedata)
// 				// .services.Nodedata.CreateNodedata(userID, inputNodedata)
// 				if err != nil {
// 					appG.ResponseError(http.StatusBadRequest, err, nil)
// 					return nil, err
// 				}

// 				node.Data = append(node.Data, *Nodedata)
// 			}
// 		}
// 		result = node
// 	}
// 	//  else {
// 	// 	fmt.Println("not found osm", bodyResponse.OsmID)
// 	// }

// 	return result, nil
// }

// func (h *HandlerV1) CreateListNode(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil || userID == "" {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return
// 	}

// 	var input []*model.Node
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	if len(input) == 0 {
// 		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
// 		return
// 	}

// 	var result []*model.Node
// 	for i := range input {
// 		// existOsmID, err := h.services.Node.FindNode(domain.RequestParams{
// 		// 	Options: domain.Options{Limit: 1},
// 		// 	Filter:  bson.D{{"osm_id", input[i].OsmID}},
// 		// })
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return
// 		// }

// 		// existLatLon := false
// 		// if len(existOsmID.Data) > 0 {
// 		// 	existLatLon = input[i].Lat == existOsmID.Data[0].Lat && input[i].Lon == existOsmID.Data[0].Lon
// 		// 	progress := 0
// 		// 	if existLatLon {
// 		// 		progress = 100
// 		// 	}

// 		// 	_, err := h.services.Ticket.CreateTicket(userID, &model.Ticket{
// 		// 		Title:       "Double osm object",
// 		// 		Description: fmt.Sprintf("[osmId]%s[/osmId]: [coords]%v,%v[/coords], [existCoords]%v,%v[/existCoords]", input[i].OsmID, input[i].Lat, input[i].Lon, existOsmID.Data[0].Lat, existOsmID.Data[0].Lon),
// 		// 		Status:      !existLatLon,
// 		// 		Progress:    progress,
// 		// 	})
// 		// 	if err != nil {
// 		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 		return
// 		// 	}
// 		// 	// fmt.Println("Double node:::", input[i].OsmID, input[i].Lat, input[i].Lon)
// 		// }
// 		// if !existLatLon {

// 		// // Get address.
// 		// pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", input[i].Lat, input[i].Lon))
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return
// 		// }
// 		// r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
// 		// r.Header.Add("User-Agent", "a127.0.0.1")

// 		// resp, err := http.DefaultClient.Do(r)
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return
// 		// }
// 		// defer resp.Body.Close()

// 		// bytes, err := io.ReadAll(resp.Body)
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return
// 		// }
// 		// var bodyResponse domain.ResponseNominatim
// 		// if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, e, nil)
// 		// 	return
// 		// }

// 		// if bodyResponse.OsmID != 0 {
// 		// 	address, err := h.services.Address.CreateAddress(userID, &domain.AddressInput{
// 		// 		OsmID:    fmt.Sprintf("%s/%d", bodyResponse.OsmType, bodyResponse.OsmID),
// 		// 		Address:  bodyResponse.Address,
// 		// 		DAddress: bodyResponse.DisplayName,
// 		// 	})
// 		// 	if err != nil {
// 		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 		return
// 		// 	}

// 		// 	input[i].OsmID = address.OsmID

// 		// 	if bodyResponse.Name == "" {
// 		// 		arrStr := strings.Split(address.DAddress, ",")
// 		// 		nameNode := ""
// 		// 		if len(arrStr) >= 2 {
// 		// 			nameNode = fmt.Sprintf("%s, %s", arrStr[1], arrStr[0])
// 		// 		} else {
// 		// 			nameNode = arrStr[0]
// 		// 		}
// 		// 		input[i].Name = strings.TrimSpace(nameNode)
// 		// 	} else {
// 		// 		input[i].Name = bodyResponse.Name
// 		// 	}

// 		// 	if ccode, ok := bodyResponse.Address["country_code"]; ok {
// 		// 		input[i].CCode = ccode.(string)
// 		// 	}
// 		// }

// 		// Node, err := h.services.Node.CreateNode(userID, input[i])
// 		// if err != nil {
// 		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 	return
// 		// }

// 		// if len(input[i].Data) > 0 {
// 		// 	for j := range input[i].Data {
// 		// 		inputNodedata := &model.NodedataInput{
// 		// 			NodeID:   Node.ID.Hex(),
// 		// 			Data:     input[i].Data[j].Data,
// 		// 			TagID:    input[i].Data[j].TagID.Hex(),
// 		// 			TagoptID: input[i].Data[j].TagoptID.Hex(),
// 		// 		}

// 		// 		Nodedata, err := h.services.Nodedata.CreateNodedata(userID, inputNodedata)
// 		// 		if err != nil {
// 		// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 		// 			return
// 		// 		}
// 		// 		Node.Data = append(Node.Data, *Nodedata)
// 		// 	}
// 		// }

// 		Node, err := h.CreateOrExistNode(c, input[i])
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return
// 		}

// 		result = append(result, Node)
// 		// }
// 	}

// 	c.JSON(http.StatusOK, result)
// }

// @Summary Node Get all Nodes
// @Security ApiKeyAuth
// @Tags Node
// @Description get all Nodes
// @ModuleID Node
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Node
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Node [get].
// func (h *HandlerV1) GetAllNode(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	params, err := utils.GetParamsFromRequest(c, model.Node{}, &h.i18n)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	nodes, err := h.services.Node.GetAllNode(params)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, nodes)
// }
