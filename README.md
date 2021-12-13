# Benchmarks

Totally **unscientific** and mostly unrealistic benchmark that
[go-faster/ch](https://github.com/go-faster/ch) project uses to understand performance.

```sql
SELECT number FROM system.numbers LIMIT 500000000
```
```
500000000 rows in set. Elapsed: 0.503 sec.
Processed 500.07 million rows,
  4.00 GB (993.26 million rows/s., 7.95 GB/s.)
```

Note: due to row-oriented design of most libraries, overhead per single row
is significantly higher, so results can be slightly surprising.

| Name                       | Protocol | Time     | RAM     |
|----------------------------|----------|----------|---------|
| clickhouse-client (C++)    | TCP      | 0.5s     | N/A     |
| **go-faster/ch**           | **TCP**  | **0.8s** | **10M** |
| clickhouse-go              | TCP      | 35s      | 182M    |
| mailru/go-clickhouse       | HTTP     | 4m13s    | 13M     |
| clickhouse-driver (Python) | TCP      | 37s      | 60M     |
| clickhouse-rs  (Rust)      | TCP      | 27s      | 182M    |
