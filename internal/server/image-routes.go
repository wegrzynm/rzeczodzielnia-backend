package server

import (
	"Rzeczodzielnia/internal/models"
	"Rzeczodzielnia/internal/utils"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	imageDir = "./images"
	maxSize  = 50 << 20 // 50MB
)

var allowedExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

func (s *Server) UploadImageHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, usr := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}
	productId := getParamsId(params)

	if err := r.ParseMultipartForm(maxSize); err != nil {
		handleError(w, http.StatusBadRequest, "Nie można sparsować formularza: "+err.Error())
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		handleError(w, http.StatusBadRequest, "Nie przesłano plików")
		return
	}

	var images []models.Image
	for _, fileHeader := range files {
		if !isAllowedExtension(fileHeader) {
			handleError(w, http.StatusBadRequest, "Niedozwolony format pliku")
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			handleError(w, http.StatusBadRequest, "Problem z otwarciem pliku")
			return
		}
		defer file.Close()

		fileName := generateUniqueFileName(fileHeader.Filename)
		filePath := getProductImagePath(productId, fileName)
		err = saveImage(filePath, file)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Błąd zapisu pliku: "+err.Error())
			return
		}
		img := models.Image{
			ProductID: productId,
			Name:      fileName,
			Path:      filePath,
			UserId:    usr.Id,
		}
		images = append(images, img)
	}
	utils.AddOrUpdateObject(&images, false)
	msg := map[string]string{"message": "Pliki zostały zapisane pomyślnie"}
	sendJSONResponse(w, http.StatusOK, msg)
}

func (s *Server) DeleteImageHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isValid, statusCode, errMsg, _ := userAuthentication(r.Header.Get("Authorization"))
	if !isValid {
		handleError(w, statusCode, errMsg)
		return
	}

	imageID := getParamsId(params)

	image := models.GetImageById(imageID)

	if image.ID == 0 {
		handleError(w, http.StatusNotFound, "Image doesn't exist")
		return
	}
	err := os.Remove(image.Path)
	if err != nil && !os.IsNotExist(err) {
		handleError(w, http.StatusInternalServerError, "Error while deleting file: "+err.Error())
		return
	}

	models.DeleteImageById(*image)

	sendJSONResponse(w, http.StatusNoContent, "Image deleted successfully")
}

func isAllowedExtension(fileHeader *multipart.FileHeader) bool {
	ext := filepath.Ext(fileHeader.Filename)
	return allowedExtensions[ext]
}

func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return uuid.New().String() + ext
}

func getProductImagePath(productID uint, fileName string) string {
	return filepath.Join(imageDir, "products", fmt.Sprint(productID), fileName)
}

func saveImage(filePath string, content io.Reader) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating filepath: %w", err)
		}
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error while saving file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, content)
	if err != nil {
		return fmt.Errorf("Error copying file : %w", err)
	}

	return nil
}
