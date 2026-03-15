package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{
		"items": h.productService.List(),
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	product, ok := h.productService.GetByID(productID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "PRODUCT_NOT_FOUND",
				"message": "Product was not found",
				"details": nil,
			},
		})
		return
	}

	c.JSON(http.StatusOK, product)
}
