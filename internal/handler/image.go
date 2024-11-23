package handler

import (
	"net/http"
	"ogimg/internal/service"
	"ogimg/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	Handler      *Handler
	imageService service.ImageService
}

func NewImageHandler(handler *Handler, imageService service.ImageService) *ImageHandler {
	return &ImageHandler{
		Handler:      handler,
		imageService: imageService,
	}
}

func (h *ImageHandler) GetOgImageByUrl(ctx *gin.Context) {
	userUrl := ctx.Query("url")
	if userUrl == "" {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "Url is required", nil)
		return
	}

	if err := h.imageService.GetOgImageByUrl(ctx, userUrl); err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
}
