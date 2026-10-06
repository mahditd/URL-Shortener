package controller

import (
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
	type Params struct {
		Url string `json:"url"`
	}
	params := Params{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	req := dto.ShortenRequest{URL: params.Url}

	res, err := h.usecase.Shorten(req)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(201, res)
}

func (h *LinkHandler) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")

	link, err := h.usecase.GetLinkByCode(code)
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": "link not found",
		})
		return
	}

	ctx.Redirect(302, link.URL)
}

// func (h *LinkHandler) Metadata(ctx *gin.Context) {
// 	code := c.Param("code")

// 	c.JSON(200, gin.H{
// 		"code": code,
// 	})
// }
