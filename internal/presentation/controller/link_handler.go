package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/application/usecase"
	domainerrors "github.com/mahditd/url-shortener/internal/domain/errors"
)

type LinkHandler struct {
	usecase *usecase.LinkUsecase
}

func NewLinkHandler(usecase *usecase.LinkUsecase) *LinkHandler {
	return &LinkHandler{
		usecase: usecase,
	}
}

func (h *LinkHandler) Shorten(ctx *gin.Context) {
	params := dto.ShortenRequest{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	res, err := h.usecase.Shorten(params)

	if err != nil {
		if errors.Is(err, domainerrors.ErrInvalidURL) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "internal server error",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *LinkHandler) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")

	link, err := h.usecase.GetLinkByCode(code)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	ctx.Redirect(http.StatusFound, link.URL)
}

func (h *LinkHandler) Metadata(ctx *gin.Context) {

	code := ctx.Param("code")

	link, err := h.usecase.GetLinkByCode(code)

	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	response := dto.MetadataResponse{
		URL:       link.URL,
		CreatedAt: link.CreatedAt,
	}

	ctx.JSON(http.StatusOK, response)

}
