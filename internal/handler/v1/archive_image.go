package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
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

func (h *HandlerV1) RegisterArchiveImage(router *gin.RouterGroup) {
	route := router.Group("/archive_image")
	route.GET("", h.findArchiveImage)
	route.DELETE("/:id", h.SetUserFromRequest, h.deleteArchiveImage)
}

func (h *HandlerV1) findArchiveImage(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.ArchiveImageInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	ArchiveImages, err := h.Services.ArchiveImage.FindArchiveImage(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, ArchiveImages)
}

func (h *HandlerV1) deleteArchiveImage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	// ArchiveImageForRemove, err := h.services.ArchiveImage.GetArchiveImage(id)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// if ArchiveImageForRemove.Service == "" {
	// 	appG.ResponseError(http.StatusBadRequest, errors.New("not found item for remove"), nil)
	// 	return
	// } else {
	// 	pathOfRemove := fmt.Sprintf("public/%s/%s", ArchiveImageForRemove.UserID.Hex(), ArchiveImageForRemove.Service)

	// 	if ArchiveImageForRemove.ServiceID != "" {
	// 		pathOfRemove = fmt.Sprintf("%s/%s", pathOfRemove, ArchiveImageForRemove.ServiceID)
	// 	}

	// 	pathRemove := fmt.Sprintf("%s/%s%s", pathOfRemove, ArchiveImageForRemove.Path, ArchiveImageForRemove.Ext)
	// 	err := os.Remove(pathRemove)
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	}

	// 	// remove srcset.
	// 	for i := range h.ArchiveImageConfig.Sizes {
	// 		dataImg := h.ArchiveImageConfig.Sizes[i]
	// 		pathRemove = fmt.Sprintf("%s/%v-%s%s", pathOfRemove, dataImg.Size, ArchiveImageForRemove.Path, ArchiveImageForRemove.Ext) // ".webp"
	// 		err = os.Remove(pathRemove)
	// 		if err != nil {
	// 			appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		}
	// 	}

	// 	// pathRemove = fmt.Sprintf("%s/xs-%s", pathOfRemove, ArchiveImageForRemove.Path)
	// 	// err = os.Remove(pathRemove)
	// 	// if err != nil {
	// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	// }
	// 	// pathRemove = fmt.Sprintf("%s/lg-%s", pathOfRemove, ArchiveImageForRemove.Path)
	// 	// err = os.Remove(pathRemove)
	// 	// if err != nil {
	// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	// }
	// }

	ArchiveImage, err := h.Services.ArchiveImage.DeleteArchiveImage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, ArchiveImage)
}
