package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type CheckoutHandler struct {
	checkoutService *service.CheckoutService
}

func NewCheckoutHandler(checkoutService *service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: checkoutService,
	}
}

func (h *CheckoutHandler) Checkout(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	order, err := h.checkoutService.Checkout(userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyCart):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "EMPTY_CART",
					"message": "Cart is empty",
					"details": nil,
				},
			})
			return
		case errors.Is(err, service.ErrProductNotFound):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "PRODUCT_NOT_FOUND",
					"message": "One or more products could not be resolved",
					"details": nil,
				},
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "CHECKOUT_FAILED",
					"message": "Checkout could not be completed",
					"details": nil,
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, order)
}
