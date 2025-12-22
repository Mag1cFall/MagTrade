package handler

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"

	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/service"
	"github.com/gin-gonic/gin"
)

type CaptchaHandler struct {
	captchaService *service.CaptchaService
}

func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{
		captchaService: service.NewCaptchaService(cache.GetClient()),
	}
}

func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		identifier = c.ClientIP()
	}

	captchaID, code, err := h.captchaService.GenerateCaptcha(c.Request.Context(), identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate captcha"})
		return
	}

	img := h.generateCaptchaImage(code)

	c.Header("X-Captcha-ID", captchaID)
	c.Header("Content-Type", "image/png")
	_ = png.Encode(c.Writer, img)
}

func (h *CaptchaHandler) CheckNeedsCaptcha(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		identifier = c.ClientIP()
	}

	needsCaptcha := h.captchaService.NeedsCaptcha(c.Request.Context(), identifier)
	isLocked := h.captchaService.IsAccountLocked(c.Request.Context(), identifier)

	c.JSON(http.StatusOK, gin.H{
		"needs_captcha": needsCaptcha,
		"is_locked":     isLocked,
	})
}

func (h *CaptchaHandler) generateCaptchaImage(code string) image.Image {
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	textColor := color.RGBA{50, 50, 150, 255}
	for i, ch := range code {
		x := 15 + i*18
		y := 28
		h.drawChar(img, x, y, byte(ch), textColor)
	}

	return img
}

func (h *CaptchaHandler) drawChar(img *image.RGBA, x, y int, ch byte, c color.Color) {
	patterns := map[byte][][]int{
		'0': {{1, 1, 1}, {1, 0, 1}, {1, 0, 1}, {1, 0, 1}, {1, 1, 1}},
		'1': {{0, 1, 0}, {1, 1, 0}, {0, 1, 0}, {0, 1, 0}, {1, 1, 1}},
		'2': {{1, 1, 1}, {0, 0, 1}, {1, 1, 1}, {1, 0, 0}, {1, 1, 1}},
		'3': {{1, 1, 1}, {0, 0, 1}, {1, 1, 1}, {0, 0, 1}, {1, 1, 1}},
		'4': {{1, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 0, 1}, {0, 0, 1}},
		'5': {{1, 1, 1}, {1, 0, 0}, {1, 1, 1}, {0, 0, 1}, {1, 1, 1}},
		'6': {{1, 1, 1}, {1, 0, 0}, {1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
		'7': {{1, 1, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}},
		'8': {{1, 1, 1}, {1, 0, 1}, {1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
		'9': {{1, 1, 1}, {1, 0, 1}, {1, 1, 1}, {0, 0, 1}, {1, 1, 1}},
	}

	pattern, ok := patterns[ch]
	if !ok {
		return
	}

	scale := 3
	for py, row := range pattern {
		for px, v := range row {
			if v == 1 {
				for sy := 0; sy < scale; sy++ {
					for sx := 0; sx < scale; sx++ {
						img.Set(x+px*scale+sx, y-len(pattern)*scale+py*scale+sy, c)
					}
				}
			}
		}
	}
}
