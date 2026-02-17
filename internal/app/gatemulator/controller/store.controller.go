package controller

import (
	"image"
	"image/jpeg"
	"log/slog"
	"net/http"
	"strings"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/service"
)

type StoreController struct {
	storeService service.StoreService
}

func NewStoreController(storeService service.StoreService) *StoreController {
	return &StoreController{storeService}
}

func (c *StoreController) Download(w http.ResponseWriter, r *http.Request) error {
	sha1 := r.PathValue("sha1")
	slog.Info("Download", slog.String("sha1", sha1))

	file, content := c.storeService.GetFile(sha1)

	if content != nil {
		if strings.HasPrefix(file.MimeType, "image/") {
			img, ok := content.(image.Image)
			if ok {
				w.Header().Set("Content-Type", file.MimeType)
				jpeg.Encode(w, img, nil)
				return nil
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
