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

type UploadHandler struct {
	uploadDir string
	baseURL   string
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

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file provided")
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		response.BadRequest(c, "file too large, max 5MB")
		return
	}

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		response.BadRequest(c, "only image files allowed")
		return
	}

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

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		response.InternalError(c, "failed to process file")
		return
	}
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))

	if _, err := file.Seek(0, 0); err != nil {
		response.InternalError(c, "failed to reset file pointer")
		return
	}

	dateDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		response.InternalError(c, "failed to create directory")
		return
	}

	filename := hashStr[:16] + ext
	filePath := filepath.Join(fullDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		url := fmt.Sprintf("%s/%s/%s", h.baseURL, dateDir, filename)
		response.Success(c, gin.H{"url": url})
		return
	}

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
