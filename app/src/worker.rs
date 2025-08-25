use serde_json::Value;
use sqlx::PgPool;
use std::env;
use std::time::Duration;
use aws_config::meta::region::RegionProviderChain;
use aws_sdk_sqs::{Client, config::Region};

#[derive(serde::Deserialize)]
struct ClassificationResponse {
    task_type: String,
}

#[derive(serde::Deserialize)]
struct ExecutionResponse {
    result: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 環境変数とクライアントのセットアップ
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let pool = PgPool::connect(&database_url).await?;

    let endpoint_url = env::var("AWS_ENDPOINT_URL").expect("AWS_ENDPOINT_URL must be set");
    let region_provider = RegionProviderChain::first_try(Region::new("us-east-1"));
    let shared_config = aws_config::from_env()
        .region(region_provider)
        .endpoint_url(endpoint_url)
        .load()
        .await;
    let sqs_client = Client::new(&shared_config);
    
    let queue_url = env::var("SQS_QUEUE_URL").expect("SQS_QUEUE_URL must be set");
    let python_agent_url = env::var("PYTHON_AGENT_URL").expect("PYTHON_AGENT_URL must be set");
    let http_client = reqwest::Client::new();

    println!("Worker started...");

    loop {
        let messages = sqs_client
            .receive_message()
            .queue_url(&queue_url)
            .max_number_of_messages(1)
            .wait_time_seconds(20)
            .send()
            .await?;

        if let Some(messages) = messages.messages {
            for message in messages {
                if let (Some(body), Some(receipt_handle)) = (message.body(), message.receipt_handle()) {
                    let payload: Value = serde_json::from_str(body)?;
                    println!("Processing task: {:?}", payload);

                    // 1. DBにタスクを登録
                    let task_id: i32 = sqlx::query_scalar!(
                        "INSERT INTO tasks (payload, status) VALUES ($1, 'pending') RETURNING id",
                        payload
                    )
                    .fetch_one(&pool)
                    .await?;

                    // 2. 分類エージェントを呼び出し
                    let classify_url = format!("{}/classify", python_agent_url);
                    let res = http_client.post(classify_url).json(&payload).send().await?;
                    let classification: ClassificationResponse = res.json().await?;
                    
                    // 3. DBを更新
                    sqlx::query!(
                        "UPDATE tasks SET status = 'classified', task_type = $1 WHERE id = $2",
                        classification.task_type,
                        task_id
                    )
                    .execute(&pool)
                    .await?;
                    println!("Task {} classified as {}", task_id, classification.task_type);


                    // 4. 専門エージェントを呼び出し
                    let execute_url = format!("{}/execute/{}", python_agent_url, classification.task_type);
                    let res = http_client.post(execute_url).json(&payload).send().await?;
                    let execution_result: ExecutionResponse = res.json().await?;

                    // 5. 最終結果をDBに保存
                    sqlx::query!(
                        "UPDATE tasks SET status = 'completed', result = $1 WHERE id = $2",
                        execution_result.result,
                        task_id
                    )
                    .execute(&pool)
                    .await?;
                    println!("Task {} completed", task_id);

                    // 6. SQSからメッセージを削除
                    sqs_client
                        .delete_message()
                        .queue_url(&queue_url)
                        .receipt_handle(receipt_handle)
                        .send()
                        .await?;
                }
            }
        } else {
            println!("No messages in queue. Waiting...");
        }
        tokio::time::sleep(Duration::from_secs(5)).await;
    }
}