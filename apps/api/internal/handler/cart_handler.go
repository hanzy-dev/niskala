package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type CartHandler struct {
	cartService *service.CartService
}

type AddCartItemRequest struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type UpdateCartItemRequest struct {
	Qty int `json:"qty"`
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	c.JSON(http.StatusOK, h.cartService.GetCart(userID))
}

func (h *CartHandler) AddCartItem(c *gin.Context) {
	var request AddCartItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "Invalid request body",
				"details": nil,
			},
		})
		return
	}

	if request.ProductID == "" || request.Qty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "product_id and qty must be valid",
				"details": nil,
			},
		})
		return
	}

	userID, _ := auth.GetUserID(c)
	c.JSON(http.StatusOK, h.cartService.AddItem(userID, request.ProductID, request.Qty))
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	var request UpdateCartItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "Invalid request body",
				"details": nil,
			},
		})
		return
	}

	userID, _ := auth.GetUserID(c)
	productID := c.Param("productId")

	c.JSON(http.StatusOK, h.cartService.UpdateItem(userID, productID, request.Qty))
}

func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	productID := c.Param("productId")

	c.JSON(http.StatusOK, h.cartService.RemoveItem(userID, productID))
}
