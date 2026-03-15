use axum::{
    routing::{get, post},
    Router,
};

use crate::handlers::{health, pricing};

pub fn build_router() -> Router {
    Router::new()
        .route("/health", get(health::get_health))
        .route("/pricing/quote", post(pricing::post_quote))
}
