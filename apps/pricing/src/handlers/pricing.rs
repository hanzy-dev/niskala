use axum::Json;

use crate::models::pricing::{
    PricingBreakdown, PricingLineItem, PricingQuoteRequest, PricingQuoteResponse,
};

pub async fn post_quote(
    Json(payload): Json<PricingQuoteRequest>,
) -> Json<PricingQuoteResponse> {
    Json(build_quote(payload))
}

pub fn build_quote(payload: PricingQuoteRequest) -> PricingQuoteResponse {
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

    PricingQuoteResponse {
        subtotal_cents,
        discount_cents,
        total_cents,
        breakdown: PricingBreakdown {
            line_items,
            coupon_code: payload.coupon_code,
            pricing_strategy: "base_price_only".to_string(),
        },
    }
}

#[cfg(test)]
mod tests {
    use super::build_quote;
    use crate::models::pricing::{PricingItem, PricingQuoteRequest};

    #[test]
    fn builds_base_price_quote_correctly() {
        let payload = PricingQuoteRequest {
            items: vec![
                PricingItem {
                    product_id: "prod_1".to_string(),
                    price_cents: 150000,
                    qty: 2,
                },
                PricingItem {
                    product_id: "prod_2".to_string(),
                    price_cents: 50000,
                    qty: 1,
                },
            ],
            coupon_code: None,
        };

        let quote = build_quote(payload);

        assert_eq!(quote.subtotal_cents, 350000);
        assert_eq!(quote.discount_cents, 0);
        assert_eq!(quote.total_cents, 350000);
        assert_eq!(quote.breakdown.line_items.len(), 2);
    }
}