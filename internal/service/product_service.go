package service

import (
	"context"

	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepository(),
	}
}

type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required,max=200"`
	Description   string  `json:"description"`
	OriginalPrice float64 `json:"original_price" binding:"required,gt=0"`
	ImageURL      string  `json:"image_url"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name" binding:"max=200"`
	Description   string  `json:"description"`
	OriginalPrice float64 `json:"original_price" binding:"omitempty,gt=0"`
	ImageURL      string  `json:"image_url"`
	Status        *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

type ProductListResponse struct {
	Products []model.Product `json:"products"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:          req.Name,
		Description:   req.Description,
		OriginalPrice: req.OriginalPrice,
		ImageURL:      req.ImageURL,
		Status:        model.ProductStatusOnShelf,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, page, pageSize int) (*ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}

	products, total, err := s.productRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *ProductService) Update(ctx context.Context, id int64, req *UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.OriginalPrice > 0 {
		product.OriginalPrice = req.OriginalPrice
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.Status != nil {
		product.Status = model.ProductStatus(*req.Status)
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}
