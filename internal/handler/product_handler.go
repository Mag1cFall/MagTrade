// 商品相關 HTTP 處理器
//
// 本檔案處理商品 CRUD 相關請求
// 列表和詳情公開訪問，建立/更新/刪除需要管理員權限
package handler

import (
	"strconv"

	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/service"
	"github.com/gin-gonic/gin"
)

// ProductHandler 商品處理器
type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

// List 查詢商品列表
// GET /api/v1/products?page=1&page_size=20
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.productService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetByID 查詢商品詳情
// GET /api/v1/products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "product not found")
		return
	}

	response.Success(c, product)
}

// Create 建立商品（管理員專用）
// POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, product)
}

// Update 更新商品（管理員專用）
// PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, product)
}

// Delete 刪除商品（軟刪除，管理員專用）
// DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
