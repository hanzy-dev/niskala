package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type AdminProductHandler struct {
	productService *service.ProductService
}

type CreateProductRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Stock       int    `json:"stock"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock"`
}

func NewAdminProductHandler(productService *service.ProductService) *AdminProductHandler {
	return &AdminProductHandler{
		productService: productService,
	}
}

func (h *AdminProductHandler) CreateProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.BadRequest(c, "INVALID_INPUT", "Body request tidak valid")
		return
	}

	if strings.TrimSpace(request.ID) == "" || strings.TrimSpace(request.Name) == "" || request.PriceCents < 0 || request.Stock < 0 {
		httpx.BadRequest(c, "INVALID_INPUT", "id, name, price_cents, dan stock harus valid")
		return
	}

	product, err := h.productService.Create(c.Request.Context(), domain.Product{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
		PriceCents:  request.PriceCents,
		Stock:       request.Stock,
		Category:    request.Category,
		ImageURL:    request.ImageURL,
		IsActive:    request.IsActive,
	})
	if err != nil {
		httpx.Internal(c, "PRODUCT_CREATE_FAILED", "Gagal membuat produk")
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *AdminProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	currentProduct, ok, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		httpx.Internal(c, "PRODUCT_GET_FAILED", "Gagal memuat produk")
		return
	}
	if !ok {
		httpx.NotFound(c, "PRODUCT_NOT_FOUND", "Produk tidak ditemukan")
		return
	}

	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.BadRequest(c, "INVALID_INPUT", "Body request tidak valid")
		return
	}

	if strings.TrimSpace(request.Name) == "" || request.PriceCents < 0 {
		httpx.BadRequest(c, "INVALID_INPUT", "name dan price_cents harus valid")
		return
	}

	updatedProduct, err := h.productService.Update(c.Request.Context(), domain.Product{
		ID:          currentProduct.ID,
		Name:        request.Name,
		Description: request.Description,
		PriceCents:  request.PriceCents,
		Stock:       currentProduct.Stock,
		Category:    request.Category,
		ImageURL:    request.ImageURL,
		IsActive:    request.IsActive,
	})
	if err != nil {
		httpx.Internal(c, "PRODUCT_UPDATE_FAILED", "Gagal memperbarui produk")
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (h *AdminProductHandler) UpdateStock(c *gin.Context) {
	productID := c.Param("id")

	var request UpdateStockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.BadRequest(c, "INVALID_INPUT", "Body request tidak valid")
		return
	}

	if request.Stock < 0 {
		httpx.BadRequest(c, "INVALID_INPUT", "stock harus bernilai 0 atau lebih")
		return
	}

	_, ok, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		httpx.Internal(c, "PRODUCT_GET_FAILED", "Gagal memuat produk")
		return
	}
	if !ok {
		httpx.NotFound(c, "PRODUCT_NOT_FOUND", "Produk tidak ditemukan")
		return
	}

	if err := h.productService.UpdateStock(c.Request.Context(), productID, request.Stock); err != nil {
		httpx.Internal(c, "PRODUCT_STOCK_UPDATE_FAILED", "Gagal memperbarui stok produk")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    productID,
		"stock": request.Stock,
	})
}
