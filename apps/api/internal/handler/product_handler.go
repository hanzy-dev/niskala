package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.productService.List(c.Request.Context())
	if err != nil {
		httpx.Internal(c, "PRODUCT_LIST_FAILED", "Failed to load products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": products,
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	product, ok, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		httpx.Internal(c, "PRODUCT_GET_FAILED", "Failed to load product")
		return
	}

	if !ok {
		httpx.NotFound(c, "PRODUCT_NOT_FOUND", "Product was not found")
		return
	}

	c.JSON(http.StatusOK, product)
}
