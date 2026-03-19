package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	orders, err := h.orderService.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		httpx.Internal(c, "ORDER_LIST_FAILED", "Gagal memuat daftar pesanan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": orders,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orderID := c.Param("id")

	order, ok, err := h.orderService.GetByUserIDAndOrderID(c.Request.Context(), userID, orderID)
	if err != nil {
		httpx.Internal(c, "ORDER_GET_FAILED", "Gagal memuat detail pesanan")
		return
	}
	if !ok {
		httpx.NotFound(c, "ORDER_NOT_FOUND", "Pesanan tidak ditemukan")
		return
	}

	c.JSON(http.StatusOK, order)
}
