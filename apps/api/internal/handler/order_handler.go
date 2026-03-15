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

	c.JSON(http.StatusOK, gin.H{
		"items": h.orderService.ListByUserID(userID),
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orderID := c.Param("id")

	order, ok := h.orderService.GetByUserIDAndOrderID(userID, orderID)
	if !ok {
		httpx.NotFound(c, "ORDER_NOT_FOUND", "Order was not found")
		return
	}

	c.JSON(http.StatusOK, order)
}
