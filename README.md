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
| **go-faster/ch**            | **TCP**  | 0.44s | 10M  | 1x      |
| clickhouse-client (C++)     | TCP      | 0.5s  | N/A  | 1.14x   |
| clickhouse-cpp (C++)        | TCP      | 0.64s | 6.7M | 1.45x   |
| vahid-sohrabloo/chconn (Go) | TCP      | 5s    | 10M  | 11x     |
| clickhouse-jdbc (Java)      | HTTP     | 10s   | 702M | 22x     |
| clickhouse-rs (Rust)        | TCP      | 27s   | 182M | 61x     |
| clickhouse-go               | TCP      | 35s   | 184M | 79x     |
| clickhouse-driver (Python)  | TCP      | 37s   | 60M  | 84x     |
| mailru/go-clickhouse        | HTTP     | 4m13s | 13M  | 575x    |

## Note about results

Benchmarks were performed on `Ryzen 9 5950x`.
Example result for ch:
```console
$ go run ./ch-bench-faster
440ms 500000000 rows 4.0 GB 9.1 GB/s
```
