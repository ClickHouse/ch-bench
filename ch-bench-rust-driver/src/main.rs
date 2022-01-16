use clickhouse_driver::prelude::*;
use std::{error::Error, time::Instant};

async fn execute(database_url: String) -> Result<(), Box<dyn Error>> {
    let pool = Pool::create(database_url.as_str())?;
    let mut conn = pool.connection().await?;
    let mut total: u64 = 0;
    let start = Instant::now();

    let mut result = conn
        .query("SELECT number FROM system.numbers_mt LIMIT 500000000")
        .await?;

    while let Some(block) = result.next().await? {
        total += block.row_count();
    }
    let elapsed = start.elapsed();
    println!("Rows: {}, elspsed: {} ms", total, elapsed.as_millis());
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    execute("tcp://localhost:9000?compression=none".to_string()).await?;

    Ok(())
}
