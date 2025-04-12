package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
)

type VMode struct {
	Quality int
	Name    string
	Size    int
	Resize  bool
	Prefix  string
	Ext     string
}

type VImages struct {
	Images []VMode
}

func UploadResizeMultipleFile(c *gin.Context, info *domain.ImageInput, nameField string, imageConfig *config.IImageConfig) ([]domain.IImagePaths, error) {
	filePaths := []domain.IImagePaths{}
	// fmt.Println("filePaths", filePaths)
	// c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(30<<20))
	form, err := c.MultipartForm()
	if err != nil {
		return filePaths, err
	}

	pathDir := fmt.Sprintf("public/%s", info.Service)
	if info.ServiceID != "" {
		pathDir = fmt.Sprintf("%s/%s", pathDir, info.ServiceID)
	}
	if _, err := os.Stat(pathDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathDir, os.ModePerm)
		if err != nil {
			// log.Println(err)
			return filePaths, err
		}
	}

	files := form.File[nameField]
	// fmt.Println("files: ", files)

	for _, file := range files {
		objImages := VImages{}
		fileExt := filepath.Ext(file.Filename)

		// add original.
		originalFileName := EncodeRus(strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename)))

		now := time.Now()
		filenameOriginal := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix())
		filePaths = append(filePaths, domain.IImagePaths{
			Path: filenameOriginal,
			Ext:  fileExt,
		})
		objImages.Images = append(objImages.Images, VMode{
			Quality: 100,
			Name:    filenameOriginal,
			Ext:     fileExt,
			Resize:  true,
			Size:    2000,
		})

		// add adaptive.
		for i := range imageConfig.Sizes {
			dataSize := imageConfig.Sizes[i]
			filenameLg := fmt.Sprintf("%v-%v-%v", dataSize.Prefix, strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-"), now.Unix())

			objImages.Images = append(objImages.Images, VMode{
				Quality: dataSize.Quality,
				Name:    filenameLg,
				Resize:  true,
				Size:    dataSize.Size,
				Ext:     fileExt,
			})
		}

		// add xs.
		// filenameXs := "xs-" + strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix())
		// objImages.Images = append(objImages.Images, VMode{
		// 	Quality: 100,
		// 	Name:    filenameXs,
		// 	Resize:  true,
		// 	Size:    30,
		// 	Ext:     fileExt,
		// })
		imaging.AutoOrientation(true)

		readerFile, _ := file.Open()
		imageFile, err := imaging.Decode(readerFile, imaging.AutoOrientation(true))
		if err != nil {
			// log.Fatal(err)
			return filePaths, err
		}
		// fmt.Println(filePaths)

		// create images.
		for i := range objImages.Images {
			dataImg := objImages.Images[i]
			imageForSave := imageFile

			// imageForSave, err := imaging.Open(imageForSave, imaging.AutoOrientation(true))

			if dataImg.Resize == true {
				imageForSave = imaging.Resize(imageForSave, dataImg.Size, 0, imaging.Lanczos) //Fill(imageForSave, 600, 600, imaging.Center, imaging.Lanczos)
				//Resize(imageForSave, dataImg.Size, 0, imaging.Lanczos)
			}

			err = imaging.Save(imageForSave, fmt.Sprintf("%s/%v%v", pathDir, dataImg.Name, dataImg.Ext), imaging.JPEGQuality(dataImg.Quality))
			if err != nil {
				return filePaths, err
			}
		}

		// encode images to webp.
		// for i := range objImages.Images {
		// 	dataImg := objImages.Images[i]
		// 	// encode webp if not original image.
		// 	if dataImg.Resize {
		// 		imagePath := fmt.Sprintf("%s/%v%v", pathDir, dataImg.Name, dataImg.Ext)
		// 		imgWebp, err := imaging.Open(imagePath, imaging.AutoOrientation(true))
		// 		if err != nil {
		// 			return filePaths, err
		// 		}
		// 		fileWebp, err := os.Create(fmt.Sprintf("%s/%v%v", pathDir, dataImg.Name, ".webp"))
		// 		if err != nil {
		// 			return filePaths, err
		// 		}
		// 		if err := webp.Encode(fileWebp, imgWebp, &webp.Options{
		// 			Lossless: false,
		// 			Quality:  float32(dataImg.Quality),
		// 			Exact:    true,
		// 		}); err != nil {
		// 			return filePaths, err
		// 		}
		// 		if err := fileWebp.Close(); err != nil {
		// 			return filePaths, err
		// 		}
		// 		err = os.Remove(imagePath)
		// 		if err != nil {
		// 			return filePaths, err
		// 		}
		// 	}
		// }

		// err = imaging.Save(imageFile, fmt.Sprintf("%s/%v", pathDir, filenameOriginal), imaging.JPEGQuality(100))
		// if err != nil {
		// 	return filePaths, err
		// }
		// src := imaging.Resize(imageFile, 1000, 0, imaging.Lanczos)
		// err = imaging.Save(src, fmt.Sprintf("%s/%v", pathDir, filenameLg), imaging.JPEGQuality(100))
		// if err != nil {
		// 	return filePaths, err
		// }
		// src_xs := imaging.Resize(imageFile, 30, 0, imaging.Lanczos)
		// err = imaging.Save(src_xs, fmt.Sprintf("%s/%v", pathDir, filenameXs), imaging.JPEGQuality(100))
		// if err != nil {
		// 	return filePaths, err
		// }

		// // webp.
		// imgWebp, err := imaging.Open(fmt.Sprintf("%s/%v", pathDir, filenameOriginal), imaging.AutoOrientation(true))
		// if err != nil {
		// 	return filePaths, err
		// }
		// filenameWebp := "webp-" + strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + ".webp"
		// fmt.Println("filenameWebp, ", filenameWebp)
		// fileWebp, err := os.Create(fmt.Sprintf("%s/%v", pathDir, filenameWebp))
		// if err != nil {
		// 	return filePaths, err
		// }
		// if err := webp.Encode(fileWebp, imgWebp, &webp.Options{
		// 	Lossless: false,
		// 	Quality:  70,
		// 	Exact:    true,
		// }); err != nil {
		// 	return filePaths, err
		// }
		// if err := fileWebp.Close(); err != nil {
		// 	return filePaths, err
		// }
	}

	// ctx.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
	return filePaths, nil
}

func UploadResizeMultipleFileForMessage(c *gin.Context, info *domain.MessageImage, nameField string, imageConfig *config.IImageConfig) ([]domain.IImagePaths, error) {
	filePaths := []domain.IImagePaths{}
	// fmt.Println("filePaths", filePaths)
	// c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(30<<20))
	form, err := c.MultipartForm()
	if err != nil {
		return filePaths, err
	}

	files := form.File[nameField]
	if len(files) == 0 {
		return filePaths, nil
	}

	pathDir := fmt.Sprintf("public/%s", info.Service)
	if info.ServiceID != "" {
		pathDir = fmt.Sprintf("%s/%s", pathDir, info.ServiceID)
	}
	if _, err := os.Stat(pathDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathDir, os.ModePerm)
		if err != nil {
			// log.Println(err)
			return filePaths, err
		}
	}

	// fmt.Println("files: ", files)

	for _, file := range files {
		// fmt.Println("filepath.Ext(file.Filename): ", filepath.Ext(file.Filename))
		imageTypes := []string{".jpg", ".jpeg", ".png", ".webp", ".ico", ".tif", ".bmp", ".gif"}
		objImages := VImages{}
		fileExt := filepath.Ext(file.Filename)

		// add original.
		originalFileName := EncodeRus(strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename)))

		now := time.Now()
		filenameOriginal := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix())
		filePaths = append(filePaths, domain.IImagePaths{
			Path: filenameOriginal,
			Ext:  fileExt,
		})

		if Contains(imageTypes, filepath.Ext(file.Filename)) {
			objImages.Images = append(objImages.Images, VMode{
				Quality: 100,
				Name:    filenameOriginal,
				Ext:     fileExt,
				Resize:  true,
				Size:    2000,
			})
			// add adaptive.
			for i := range imageConfig.Sizes {
				dataSize := imageConfig.Sizes[i]
				filenameLg := fmt.Sprintf("%v-%v-%v", dataSize.Prefix, strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-"), now.Unix())

				objImages.Images = append(objImages.Images, VMode{
					Quality: dataSize.Quality,
					Name:    filenameLg,
					Resize:  true,
					Size:    dataSize.Size,
					Ext:     fileExt,
				})
			}

			imaging.AutoOrientation(true)

			readerFile, _ := file.Open()
			imageFile, err := imaging.Decode(readerFile, imaging.AutoOrientation(true))
			if err != nil {
				// log.Fatal(err)
				return filePaths, err
			}
			// fmt.Println(filePaths)

			// create images.
			for i := range objImages.Images {
				dataImg := objImages.Images[i]
				imageForSave := imageFile

				// imageForSave, err := imaging.Open(imageForSave, imaging.AutoOrientation(true))

				if dataImg.Resize == true {
					imageForSave = imaging.Resize(imageForSave, dataImg.Size, 0, imaging.Lanczos) //Fill(imageForSave, 600, 600, imaging.Center, imaging.Lanczos)
					//Resize(imageForSave, dataImg.Size, 0, imaging.Lanczos)
				}

				err = imaging.Save(imageForSave, fmt.Sprintf("%s/%v%v", pathDir, dataImg.Name, dataImg.Ext), imaging.JPEGQuality(dataImg.Quality))
				if err != nil {
					return filePaths, err
				}
			}
		} else {
			readerFile, _ := file.Open()
			// Start reading multi-part file under id "fileupload"
			// f, fh, err := c.Request.FormFile("images")
			// if err != nil {
			// 	if err == http.ErrMissingFile {
			// 		return filePaths, err
			// 	} else {
			// 		return filePaths, err
			// 	}

			// 	return filePaths, err
			// }
			// defer f.Close()

			destDir, err := os.MkdirTemp("", "")
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)

				return filePaths, err
			}

			// Remove destDir in case of error
			defer func() {
				if err != nil {
					if remErr := os.RemoveAll(destDir); remErr != nil {
						// Log some kind of warning probably?
					}
				}
			}()

			// Create a file in our temporary directory with the same name
			// as the uploaded file
			destFile, err := os.Create(fmt.Sprintf("%s/%v%v", pathDir, filenameOriginal, fileExt)) //path.Join(destDir, fh.Filename)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)

				return filePaths, err
			}
			defer destFile.Close()
			// fmt.Println("destFile: ", destFile.Name(), fmt.Sprintf("%s/%v", pathDir, fh.Filename))

			// Write contents of uploaded file to destFile
			if _, err = io.Copy(destFile, readerFile); err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)

				return filePaths, err
			}
			// filePaths = append(filePaths, domain.IImagePaths{
			// 	Path: destFile.Name(),
			// 	Ext:  filepath.Ext(file.Filename),
			// })
		}

	}

	// ctx.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
	return filePaths, nil
}
