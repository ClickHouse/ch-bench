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

| Name                        | Protocol | Time  | RAM  | Speedup |
|-----------------------------|----------|-------|------|---------|
| clickhouse-client (C++)     | TCP      | 0.5s  | N/A  | 0.65x   |
| clickhouse-cpp (C++)        | TCP      | 0.64s | 6.7M | 0.91x   |
| **go-faster/ch**            | **TCP**  | 0.7s  | 10M  | 1x      |
| vahid-sohrabloo/chconn (Go) | TCP      | 5s    | 10M  | 7x      |
| clickhouse-rs (Rust)        | TCP      | 27s   | 182M | 38x     |
| clickhouse-jdbc (Java)      | HTTP     | 27s   | 271M | 38x     |
| clickhouse-go               | TCP      | 35s   | 184M | 50x     |
| clickhouse-driver (Python)  | TCP      | 37s   | 60M  | 52x     |
| mailru/go-clickhouse        | HTTP     | 4m13s | 13M  | 360x    |
