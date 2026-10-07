package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/application/usecase"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *LinkHandler) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")

	link, err := h.usecase.GetLinkByCode(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "link not found",
		})
		return
	}

	ctx.Redirect(http.StatusFound, link.URL)
}

func (h *LinkHandler) Metadata(ctx *gin.Context) {

	code := ctx.Param("code")

	link, err := h.usecase.GetLinkByCode(code)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "link not found",
		})
		return
	}

	response := dto.MetadataResponse{
		URL:       link.URL,
		CreatedAt: link.CreatedAt,
	}

	ctx.JSON(http.StatusOK, response)

}
