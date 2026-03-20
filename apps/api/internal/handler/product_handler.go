package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type ProductHandler struct {
	productService      *service.ProductService
	productCacheService *service.ProductCacheService
}

func NewProductHandler(
	productService *service.ProductService,
	productCacheService *service.ProductCacheService,
) *ProductHandler {
	return &ProductHandler{
		productService:      productService,
		productCacheService: productCacheService,
	}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	if products, ok := h.productCacheService.GetProductList(c.Request.Context()); ok {
		c.JSON(http.StatusOK, gin.H{"items": products})
		return
	}

	products, err := h.productService.List(c.Request.Context())
	if err != nil {
		httpx.Internal(c, "PRODUCT_LIST_FAILED", "Failed to load products")
		return
	}

	h.productCacheService.SetProductList(c.Request.Context(), products)

	c.JSON(http.StatusOK, gin.H{
		"items": products,
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	if product, ok := h.productCacheService.GetProduct(c.Request.Context(), productID); ok {
		c.JSON(http.StatusOK, product)
		return
	}

	product, ok, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		httpx.Internal(c, "PRODUCT_GET_FAILED", "Failed to load product")
		return
	}

	if !ok {
		httpx.NotFound(c, "PRODUCT_NOT_FOUND", "Product was not found")
		return
	}

	h.productCacheService.SetProduct(c.Request.Context(), product)

	c.JSON(http.StatusOK, product)
}
