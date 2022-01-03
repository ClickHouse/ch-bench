use std::{error::Error};
use clickhouse_rs::{Pool};
use futures_util::{TryStreamExt};
use futures::future;

async fn execute(database_url: String) -> Result<(), Box<dyn Error>> {
    let pool = Pool::new(database_url);

    let mut client = pool.get_handle().await.unwrap();
    let mut total: u64 = 0;

    client.query("SELECT number FROM system.numbers_mt LIMIT 500000000")
        .stream_blocks()
        .try_for_each(|block| {
            total += block.row_count() as u64;
            future::ready(Ok(()))
        }).await?;

    println!("Rows: {}", total);

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    execute("tcp://localhost:9000".to_string()).await?;

    Ok(())
}
