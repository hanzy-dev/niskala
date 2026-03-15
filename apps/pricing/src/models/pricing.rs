use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct PricingQuoteRequest {
    pub items: Vec<PricingItem>,
    pub coupon_code: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct PricingItem {
    pub product_id: String,
    pub price_cents: i64,
    pub qty: i64,
}

#[derive(Debug, Serialize)]
pub struct PricingQuoteResponse {
    pub subtotal_cents: i64,
    pub discount_cents: i64,
    pub total_cents: i64,
    pub breakdown: PricingBreakdown,
}

#[derive(Debug, Serialize)]
pub struct PricingBreakdown {
    pub line_items: Vec<PricingLineItem>,
    pub coupon_code: Option<String>,
    pub pricing_strategy: String,
}
#[derive(Debug, Serialize)]
pub struct PricingLineItem {
    pub product_id: String,
    pub price_cents: i64,
    pub qty: i64,
    pub line_subtotal_cents: i64,
}
