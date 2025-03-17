package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/middleware"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"github.com/mikalai2006/kingwood-api/pkg/app"
)

// func init() {
// 	if _, err := os.Stat("public/single"); errors.Is(err, os.ErrNotExist) {
// 		err := os.MkdirAll("public/single", os.ModePerm)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// 	if _, err := os.Stat("public/multiple"); errors.Is(err, os.ErrNotExist) {
// 		err := os.MkdirAll("public/multiple", os.ModePerm)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

func (h *HandlerV1) RegisterImage(router *gin.RouterGroup) {
	route := router.Group("/image")
	route.POST("", h.SetUserFromRequest, h.createImage)
	route.GET("", h.findImage)
	route.GET("/:id", h.getImage)
	route.GET("/:id/dir", h.SetUserFromRequest, h.getImageDirs)
	route.DELETE("/:id", h.SetUserFromRequest, h.deleteImage)
}

func (h *HandlerV1) getImage(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	image, err := h.Services.Image.GetImage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, image)
}

func (h *HandlerV1) getImageDirs(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	image, err := h.Services.Image.GetImageDirs(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, image)
}

func (h *HandlerV1) findImage(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.ImageInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	images, err := h.Services.Image.FindImage(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, images)
}

func (h *HandlerV1) createImage(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input = &domain.ImageInput{}
	if er := c.Bind(input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println("input", input)
	input.UserID = userID
	// var image domain.Image

	paths, err := utils.UploadResizeMultipleFile(c, input, "images", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}

	var result []domain.Image
	for i := range paths {
		input.Path = paths[i].Path
		input.Ext = paths[i].Ext
		image, err := h.Services.Image.CreateImage(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, image)
	}
	c.JSON(http.StatusOK, result)
}

func (h *HandlerV1) deleteImage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	// imageForRemove, err := h.services.Image.GetImage(id)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// if imageForRemove.Service == "" {
	// 	appG.ResponseError(http.StatusBadRequest, errors.New("not found item for remove"), nil)
	// 	return
	// } else {
	// 	pathOfRemove := fmt.Sprintf("public/%s/%s", imageForRemove.UserID.Hex(), imageForRemove.Service)

	// 	if imageForRemove.ServiceID != "" {
	// 		pathOfRemove = fmt.Sprintf("%s/%s", pathOfRemove, imageForRemove.ServiceID)
	// 	}

	// 	pathRemove := fmt.Sprintf("%s/%s%s", pathOfRemove, imageForRemove.Path, imageForRemove.Ext)
	// 	err := os.Remove(pathRemove)
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	}

	// 	// remove srcset.
	// 	for i := range h.imageConfig.Sizes {
	// 		dataImg := h.imageConfig.Sizes[i]
	// 		pathRemove = fmt.Sprintf("%s/%v-%s%s", pathOfRemove, dataImg.Size, imageForRemove.Path, imageForRemove.Ext) // ".webp"
	// 		err = os.Remove(pathRemove)
	// 		if err != nil {
	// 			appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		}
	// 	}

	// 	// pathRemove = fmt.Sprintf("%s/xs-%s", pathOfRemove, imageForRemove.Path)
	// 	// err = os.Remove(pathRemove)
	// 	// if err != nil {
	// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	// }
	// 	// pathRemove = fmt.Sprintf("%s/lg-%s", pathOfRemove, imageForRemove.Path)
	// 	// err = os.Remove(pathRemove)
	// 	// if err != nil {
	// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	// }
	// }

	image, err := h.Services.Image.DeleteImage(userID, id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, image)
}
