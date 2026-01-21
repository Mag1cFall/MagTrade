// 商品模型單元測試
//
// 測試覆蓋：
// - Product.TableName: GORM 資料表名稱
// - ProductStatus 常量: 上架/下架狀態值驗證
// - Product 欄位: 模型欄位賦值與讀取
package model

import (
	"testing"
)

func TestProduct_TableName(t *testing.T) {
	p := Product{}
	if got := p.TableName(); got != "products" {
		t.Errorf("Product.TableName() = %v, want %v", got, "products")
	}
}

func TestProductStatus_Constants(t *testing.T) {
	tests := []struct {
		name   string
		status ProductStatus
		want   int8
	}{
		{"off shelf", ProductStatusOffShelf, 0},
		{"on shelf", ProductStatusOnShelf, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int8(tt.status) != tt.want {
				t.Errorf("ProductStatus %s = %v, want %v", tt.name, int8(tt.status), tt.want)
			}
		})
	}
}

func TestProduct_Fields(t *testing.T) {
	p := Product{
		ID:            1,
		Name:          "Test Product",
		Description:   "Test Description",
		OriginalPrice: 99.99,
		ImageURL:      "https://example.com/image.jpg",
		Status:        ProductStatusOnShelf,
	}

	if p.ID != 1 {
		t.Errorf("Product.ID = %v, want %v", p.ID, 1)
	}
	if p.Name != "Test Product" {
		t.Errorf("Product.Name = %v, want %v", p.Name, "Test Product")
	}
	if p.Description != "Test Description" {
		t.Errorf("Product.Description = %v, want %v", p.Description, "Test Description")
	}
	if p.OriginalPrice != 99.99 {
		t.Errorf("Product.OriginalPrice = %v, want %v", p.OriginalPrice, 99.99)
	}
	if p.ImageURL != "https://example.com/image.jpg" {
		t.Errorf("Product.ImageURL = %v, want %v", p.ImageURL, "https://example.com/image.jpg")
	}
	if p.Status != ProductStatusOnShelf {
		t.Errorf("Product.Status = %v, want %v", p.Status, ProductStatusOnShelf)
	}
}
