// 檔案上傳 HTTP 處理器
//
// 本檔案處理圖片上傳功能
// 支援重複檔案去重（MD5 檢測）
// 按日期分目錄儲存：./uploads/2024/12/xxx.jpg
package handler

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// UploadHandler 上傳處理器
type UploadHandler struct {
	uploadDir string // 上傳目錄
	baseURL   string // 訪問 URL 前綴
}

func NewUploadHandler() *UploadHandler {
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}
	return &UploadHandler{
		uploadDir: uploadDir,
		baseURL:   "/uploads",
	}
}

// Upload 上傳圖片
// POST /api/v1/upload
// 表單欄位：file（檔案）
func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file provided")
		return
	}
	defer file.Close()

	// 檢查檔案大小（最大 5MB）
	if header.Size > 5*1024*1024 {
		response.BadRequest(c, "file too large, max 5MB")
		return
	}

	// 只允許圖片
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		response.BadRequest(c, "only image files allowed")
		return
	}

	// 確定副檔名
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	// 計算檔案 MD5（用於去重）
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		response.InternalError(c, "failed to process file")
		return
	}
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))

	// 重置檔案指標
	if _, err := file.Seek(0, 0); err != nil {
		response.InternalError(c, "failed to reset file pointer")
		return
	}

	// 按日期建立目錄
	dateDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		response.InternalError(c, "failed to create directory")
		return
	}

	// 使用 MD5 前 16 位作為檔名
	filename := hashStr[:16] + ext
	filePath := filepath.Join(fullDir, filename)

	// 如果檔案已存在（重複上傳），直接返回 URL
	if _, err := os.Stat(filePath); err == nil {
		url := fmt.Sprintf("%s/%s/%s", h.baseURL, dateDir, filename)
		response.Success(c, gin.H{"url": url})
		return
	}

	// 儲存檔案
	dst, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.InternalError(c, "failed to write file")
		return
	}

	url := fmt.Sprintf("%s/%s/%s", h.baseURL, dateDir, filename)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"url": url},
	})
}
