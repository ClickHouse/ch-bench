use serde::Deserialize;

use clickhouse::{error::Result, Client, Compression, Row};

#[derive(Row, Deserialize)]
struct Data {
    no: u64,
}

#[tokio::main]
async fn main() -> Result<()> {
    let client = Client::default()
        .with_compression(Compression::None)
        .with_url("http://localhost:8123");

    let mut cursor = client
        .query("SELECT number FROM system.numbers_mt LIMIT 500000000")
        .fetch::<Data>()?;

    while let Some(_row) = cursor.next().await? {}

    Ok(())
}
