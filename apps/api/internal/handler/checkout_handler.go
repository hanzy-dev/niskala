package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
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
	idempotencyKey := c.GetHeader("Idempotency-Key")

	order, err := h.checkoutService.Checkout(c.Request.Context(), userID, idempotencyKey)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingIdempotencyKey):
			httpx.BadRequest(c, "MISSING_IDEMPOTENCY_KEY", "Idempotency-Key header is required")
			return
		case errors.Is(err, service.ErrIdempotencyInProgress):
			httpx.JSONError(c, http.StatusConflict, "IDEMPOTENCY_IN_PROGRESS", "Checkout is already being processed for this key", nil)
			return
		case errors.Is(err, service.ErrEmptyCart):
			httpx.BadRequest(c, "EMPTY_CART", "Cart is empty")
			return
		case errors.Is(err, service.ErrProductNotFound):
			httpx.BadRequest(c, "PRODUCT_NOT_FOUND", "One or more products could not be resolved")
			return
		case errors.Is(err, service.ErrInsufficientStock):
			httpx.BadRequest(c, "INSUFFICIENT_STOCK", "Requested quantity exceeds available stock")
			return
		default:
			httpx.Internal(c, "CHECKOUT_FAILED", "Checkout could not be completed")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"meta": gin.H{
			"correlation_id":        httpx.GetCorrelationID(c),
			"pricing_fallback_used": order.PricingFallbackUsed,
		},
	})
}
