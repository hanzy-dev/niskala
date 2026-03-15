mod handlers {
    pub mod health;
    pub mod pricing;
}

mod models;
mod router;

use std::net::SocketAddr;

use tokio::net::TcpListener;

#[tokio::main]
async fn main() {
    let port = std::env::var("PORT").unwrap_or_else(|_| "8081".to_string());
    let address: SocketAddr = format!("127.0.0.1:{port}")
        .parse()
        .expect("invalid bind address");

    let listener = TcpListener::bind(address)
        .await
        .expect("failed to bind pricing server");

    println!("starting pricing service on http://{address}");

    axum::serve(listener, router::build_router())
        .await
        .expect("pricing server failed");
}