use axum::Json;

use crate::models::pricing::{
    PricingBreakdown, PricingLineItem, PricingQuoteRequest, PricingQuoteResponse,
};

pub async fn post_quote(
    Json(payload): Json<PricingQuoteRequest>,
) -> Json<PricingQuoteResponse> {
    let mut subtotal_cents = 0_i64;
    let mut line_items = Vec::with_capacity(payload.items.len());

    for item in payload.items {
        let line_subtotal_cents = item.price_cents * item.qty;
        subtotal_cents += line_subtotal_cents;

        line_items.push(PricingLineItem {
            product_id: item.product_id,
            price_cents: item.price_cents,
            qty: item.qty,
            line_subtotal_cents,
        });
    }

    let discount_cents = 0_i64;
    let total_cents = subtotal_cents - discount_cents;

    Json(PricingQuoteResponse {
        subtotal_cents,
        discount_cents,
        total_cents,
        breakdown: PricingBreakdown {
            line_items,
            coupon_code: payload.coupon_code,
            pricing_strategy: "base_price_only".to_string(),
        },
    })
}