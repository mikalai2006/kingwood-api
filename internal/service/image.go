package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
	"github.com/mikalai2006/kingwood-api/internal/utils"
)

type ImageService struct {
	repo        repository.Image
	ImageConfig config.IImageConfig
	Services    *Services
}

func NewImageService(repo repository.Image, imageConfig config.IImageConfig) *ImageService {
	return &ImageService{repo: repo, ImageConfig: imageConfig}
}

func (s *ImageService) FindImage(params *domain.ImageFilter) (domain.Response[domain.Image], error) {
	return s.repo.FindImage(params)
}

func (s *ImageService) GetImage(id string) (domain.Image, error) {
	return s.repo.GetImage(id)
}

func (s *ImageService) GetImageDirs(id string) ([]interface{}, error) {
	return s.repo.GetImageDirs(id)
}
func (s *ImageService) CreateImage(userID string, image *domain.ImageInput) (domain.Image, error) {
	var result domain.Image

	if image.Service == "user" {
		// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
		// if err != nil {
		// 	return result, err
		// }

		existImage, err := s.repo.FindImage(&domain.ImageFilter{UserId: []string{userID}, Service: []string{image.Service}, ServiceId: []string{image.ServiceID}})
		// domain.RequestParams{Filter: bson.D{
		// 	{"userId", userIDPrimitive},
		// 	{"service", image.Service},
		// 	{"serviceId", image.ServiceID},
		// }}
		if err != nil {
			return result, err
		}

		if len(existImage.Data) > 0 {
			for i, _ := range existImage.Data {
				_, _ = s.DeleteImage(userID, existImage.Data[i].ID.Hex(), false)
				// if err != nil {
				// 	return result, err
				// }
			}
		}

	}
	result, err := s.repo.CreateImage(userID, image)

	return result, err
}

func (s *ImageService) DeleteImage(userID string, id string, createArchive bool) (domain.Image, error) {
	// result := domain.Image{}
	// imageForRemove, err := s.GetImage(id)
	// if err != nil {
	// 	return result, err
	// }
	// if imageForRemove.Service == "" {
	// 	return result, errors.New("not found item for remove")
	// } else {
	// 	pathOfRemove := fmt.Sprintf("public/%s", imageForRemove.Service)

	// 	if imageForRemove.ServiceID != "" {
	// 		pathOfRemove = fmt.Sprintf("%s/%s", pathOfRemove, imageForRemove.ServiceID)
	// 	}

	// 	pathRemove := fmt.Sprintf("%s/%s%s", pathOfRemove, imageForRemove.Path, imageForRemove.Ext)
	// 	os.Remove(pathRemove)
	// 	// if err != nil {
	// 	// 	return result, err
	// 	// }

	// 	// remove srcset.
	// 	for i := range s.imageConfig.Sizes {
	// 		dataImg := s.imageConfig.Sizes[i]
	// 		pathRemove = fmt.Sprintf("%s/%v-%s%s", pathOfRemove, dataImg.Prefix, imageForRemove.Path, imageForRemove.Ext) // ".webp"
	// 		// fmt.Println("pathRemove2=", pathRemove)
	// 		os.Remove(pathRemove)
	// 		// if err != nil {
	// 		// 	return result, err
	// 		// }
	// 	}

	// 	isEmpty, err := utils.IsEmptyDir(pathOfRemove)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	if isEmpty {
	// 		err = os.Remove(pathOfRemove)
	// 		if err != nil {
	// 			return result, err
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
	result, err := s.repo.DeleteImage(id)
	if err != nil {
		return result, err
	}

	if !result.ID.IsZero() {
		if createArchive {

			// add to archive.
			_, err = s.Services.CreateArchiveImage(userID, &result)
		} else {
			imageForRemove := result

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
		}
	}

	return result, err //s.repo.DeleteImage(id)
}
