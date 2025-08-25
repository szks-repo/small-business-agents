use axum::{routing::post, Json, Router};
use aws_config::meta::region::RegionProviderChain;
use aws_sdk_sqs::{Client, config::Region};
use serde_json::{json, Value};
use std::env;

async fn webhook_handler(
    Json(payload): Json<Value>,
    sqs_client: axum::Extension<Client>,
) -> Json<Value> {
    let queue_url = env::var("SQS_QUEUE_URL").expect("SQS_QUEUE_URL must be set");

    let msg_body = serde_json::to_string(&payload).unwrap();

    match sqs_client
        .send_message()
        .queue_url(queue_url)
        .message_body(msg_body)
        .send()
        .await
    {
        Ok(_) => Json(json!({"status": "accepted"})),
        Err(e) => {
            eprintln!("Failed to send message to SQS: {:?}", e);
            Json(json!({"status": "error"}))
        }
    }
}

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok();
    
    let endpoint_url = env::var("AWS_ENDPOINT_URL").expect("AWS_ENDPOINT_URL must be set");
    let region_provider = RegionProviderChain::first_try(Region::new("us-east-1"));
    let shared_config = aws_config::from_env()
        .region(region_provider)
        .endpoint_url(endpoint_url)
        .load()
        .await;
    let sqs_client = Client::new(&shared_config);

    let app = Router::new()
        .route("/webhook", post(webhook_handler))
        .layer(axum::Extension(sqs_client));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    println!("Listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}
