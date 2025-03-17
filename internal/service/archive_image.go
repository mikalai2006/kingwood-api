package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ArchiveImageService struct {
	repo        repository.ArchiveImage
	ImageConfig config.IImageConfig
}

func NewArchiveImageService(repo repository.ArchiveImage, ImageConfig config.IImageConfig) *ArchiveImageService {
	return &ArchiveImageService{repo: repo, ImageConfig: ImageConfig}
}

func (s *ArchiveImageService) FindArchiveImage(params domain.RequestParams) (domain.Response[domain.ArchiveImage], error) {
	return s.repo.FindArchiveImage(params)
}

func (s *ArchiveImageService) CreateArchiveImage(userID string, data *domain.Image) (domain.ArchiveImage, error) {
	var result domain.ArchiveImage

	result, err := s.repo.CreateArchiveImage(userID, data)

	return result, err
}

func (s *ArchiveImageService) DeleteArchiveImage(id string) (domain.ArchiveImage, error) {
	result := domain.ArchiveImage{}
	ArchiveImagesForRemove, err := s.FindArchiveImage(domain.RequestParams{Filter: bson.D{{"_id", id}}})
	if err != nil {
		return result, err
	}
	var imageForRemove domain.ArchiveImage
	if len(ArchiveImagesForRemove.Data) > 0 {
		imageForRemove = ArchiveImagesForRemove.Data[0]
	}
	if imageForRemove.Service == "" {
		return result, errors.New("not found item for remove")
	} else {
		pathOfRemove := fmt.Sprintf("public/%s", imageForRemove.Service)

		if imageForRemove.ServiceID != "" {
			pathOfRemove = fmt.Sprintf("%s/%s", pathOfRemove, imageForRemove.ServiceID)
		}

		pathRemove := fmt.Sprintf("%s/%s%s", pathOfRemove, imageForRemove.Path, imageForRemove.Ext)
		os.Remove(pathRemove)
		// if err != nil {
		// 	return result, err
		// }

		// remove srcset.
		for i := range s.ImageConfig.Sizes {
			dataImg := s.ImageConfig.Sizes[i]
			pathRemove = fmt.Sprintf("%s/%v-%s%s", pathOfRemove, dataImg.Prefix, imageForRemove.Path, imageForRemove.Ext) // ".webp"
			// fmt.Println("pathRemove2=", pathRemove)
			os.Remove(pathRemove)
			// if err != nil {
			// 	return result, err
			// }
		}

		isEmpty, err := utils.IsEmptyDir(pathOfRemove)
		if err != nil {
			return result, err
		}
		if isEmpty {
			err = os.Remove(pathOfRemove)
			if err != nil {
				return result, err
			}
		}

		// pathRemove = fmt.Sprintf("%s/xs-%s", pathOfRemove, ArchiveImageForRemove.Path)
		// err = os.Remove(pathRemove)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// }
		// pathRemove = fmt.Sprintf("%s/lg-%s", pathOfRemove, ArchiveImageForRemove.Path)
		// err = os.Remove(pathRemove)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// }
	}

	return s.repo.DeleteArchiveImage(id)
}
