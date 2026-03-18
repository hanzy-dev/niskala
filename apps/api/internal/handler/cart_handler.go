package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
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

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		httpx.Internal(c, "CART_GET_FAILED", "Gagal memuat keranjang")
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddCartItem(c *gin.Context) {
	var request AddCartItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.BadRequest(c, "INVALID_INPUT", "Body request tidak valid")
		return
	}

	if request.ProductID == "" || request.Qty <= 0 {
		httpx.BadRequest(c, "INVALID_INPUT", "product_id dan qty harus valid")
		return
	}

	userID, _ := auth.GetUserID(c)

	cart, err := h.cartService.AddItem(c.Request.Context(), userID, request.ProductID, request.Qty)
	if err != nil {
		httpx.Internal(c, "CART_ADD_FAILED", "Gagal menambahkan item ke keranjang")
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	var request UpdateCartItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.BadRequest(c, "INVALID_INPUT", "Body request tidak valid")
		return
	}

	userID, _ := auth.GetUserID(c)
	productID := c.Param("productId")

	cart, err := h.cartService.UpdateItem(c.Request.Context(), userID, productID, request.Qty)
	if err != nil {
		httpx.Internal(c, "CART_UPDATE_FAILED", "Gagal memperbarui item keranjang")
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	productID := c.Param("productId")

	cart, err := h.cartService.RemoveItem(c.Request.Context(), userID, productID)
	if err != nil {
		httpx.Internal(c, "CART_DELETE_FAILED", "Gagal menghapus item keranjang")
		return
	}

	c.JSON(http.StatusOK, cart)
}
