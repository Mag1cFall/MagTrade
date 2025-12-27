// 驗證碼 HTTP 處理器
//
// 本檔案處理圖形驗證碼生成和驗證
// 支援登入失敗多次後強制驗證碼、帳號鎖定檢測
// 使用純 Go 繪製簡易數字驗證碼圖片
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

// CaptchaHandler 驗證碼處理器
type CaptchaHandler struct {
	captchaService *service.CaptchaService
}

func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{
		captchaService: service.NewCaptchaService(cache.GetClient()),
	}
}

// GetCaptcha 取得驗證碼圖片
// GET /api/v1/captcha?identifier=xxx
// 返回 PNG 圖片，X-Captcha-ID Header 包含驗證碼 ID
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

	img := h.generateCaptchaImage(code) // 生成驗證碼圖片

	c.Header("X-Captcha-ID", captchaID)
	c.Header("Content-Type", "image/png")
	_ = png.Encode(c.Writer, img)
}

// CheckNeedsCaptcha 檢查是否需要驗證碼
// GET /api/v1/captcha/check?identifier=xxx
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

// generateCaptchaImage 生成驗證碼圖片（純 Go 實現）
func (h *CaptchaHandler) generateCaptchaImage(code string) image.Image {
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景色
	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 繪製數字
	textColor := color.RGBA{50, 50, 150, 255}
	for i, ch := range code {
		x := 15 + i*18
		y := 28
		h.drawChar(img, x, y, byte(ch), textColor)
	}

	return img
}

// drawChar 繪製單個字符（使用點陣圖案）
func (h *CaptchaHandler) drawChar(img *image.RGBA, x, y int, ch byte, c color.Color) {
	// 數字 0-9 的 3x5 點陣圖案
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

	scale := 3 // 放大倍數
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
