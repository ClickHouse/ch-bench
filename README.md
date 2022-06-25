# Benchmarks

Totally **unscientific** and mostly unrealistic benchmark that
[go-faster/ch][faster] project uses to understand performance.

The main goal is to measure minimal **client overhead** (CPU, RAM) to read data,
i.e. data blocks deserialization and transfer.

Please see [Notes](#Notes) for more details about results.

```sql
SELECT number FROM system.numbers_mt LIMIT 500000000
```
```
500000000 rows in set. Elapsed: 0.503 sec.
Processed 500.07 million rows,
  4.00 GB (993.26 million rows/s., 7.95 GB/s.)
```

Note: due to row-oriented design of most libraries, overhead per single row
is significantly higher, so results can be slightly surprising.

| Name                                       | Time  | RAM  | Ratio |
|--------------------------------------------|-------|------|-------|
| **[ClickHouse/ch-go][ch]** (Go)            | 401ms | 9M   | ~1x   |
| [clickhouse-client][client] (C++)          | 387ms | 91M  | ~1x   |
| [vahid-sohrabloo/chconn][vahid] (Go)       | 472ms | 9M   | ~1x   |
| [clickhouse-cpp][cpp] (C++)                | 516ms | 6.9M | 1.47x |
| [clickhouse_driver][rs] (Rust)             | 614ms | 9M   | 1.72x |
| [uptrace][uptrace] (Go)                    | 1.2s  | 8M   | 3x    |
| [clickhouse-go][go] (Go)                   | 5.4s  | 21M  | 13.5x |
| [curl][curl] (C, HTTP)                     | 5.4s  | 21M  | 13.5x |
| [clickhouse-client][java] (Java, HTTP)     | 6.4s  | 121M | 16x   |
| [clickhouse-jdbc][jdbc] (Java, HTTP)       | 7.2s  | 120M | 18x   |
| [loyd/clickhouse.rs][rs-http] (Rust, HTTP) | 10s   | 7.2M | 28x   |
| [clickhouse-driver][py] (Python)           | 37s   | 60M  | 106x  |
| [mailru/go-clickhouse][mail] (Go, HTTP)    | 4m13s | 13M  | 729x  |

[client]:  https://clickhouse.com/docs/en/interfaces/cli/ "Native command-line client (Official)"
[ch]:      https://github.com/ClickHouse/ch-go "ClickHouse/ch-go"
[rs]:      https://github.com/datafuse-extras/clickhouse_driver "datafuse-extras/clickhouse_driver"
[rs-http]: https://github.com/loyd/clickhouse.rs "A typed client for ClickHouse (HTTP)"
[cpp]:     https://github.com/ClickHouse/clickhouse-cpp "C++ client library for ClickHouse (Official)"
[curl]:    https://github.com/curl/curl "A command-line tool for transferring data specified with URL syntax"
[vahid]:   https://github.com/vahid-sohrabloo/chconn "Low-level ClickHouse database driver for Golang"
[java]:    https://github.com/ClickHouse/clickhouse-jdbc/tree/develop/clickhouse-client "Java client for ClickHouse (Official)"
[jdbc]:    https://github.com/ClickHouse/clickhouse-jdbc/tree/develop/clickhouse-jdbc "JDBC driver for ClickHouse (Official)"
[py]:      https://github.com/mymarilyn/clickhouse-driver
[go]:      https://github.com/ClickHouse/clickhouse-go "Golang driver for ClickHouse (Official)"
[mail]:    https://github.com/mailru/go-clickhouse "Golang SQL database driver (HTTP, TSV format)"
[uptrace]: https://github.com/uptrace/go-clickhouse "ClickHouse client for Go 1.18+ (Uptrace)"

See [RESULTS.md](./RESULTS.md) and [RESULTS.slow.md](./RESULTS.slow.md).

<sub>
Keeping `go-faster/ch`, `clickhouse-client` and `vahid-sohrabloo/chconn` to `~1x`, they are mostly equal.
</sub>

## Notes

### C++

| Command                  |      Mean [ms] | Min [ms] | Max [ms] |    Relative |
|:-------------------------|---------------:|---------:|---------:|------------:|
| `go-faster`              |   598.8 ± 92.2 |    356.9 |    792.8 | 1.07 ± 0.33 |
| `clickhouse-client`      |  561.9 ± 149.5 |    387.8 |   1114.2 |        1.00 |
| `clickhouse-cpp`         |   574.4 ± 35.9 |    523.3 |    707.4 | 1.02 ± 0.28 |


We are selecting **best** results, however C++ client has lower dispersion.

# Maximum possible speed

I've measured my localhost performance using `iperf3`, getting 10 GiB/s,
this correlates with top results.

For example, one of [go-faster/ch][faster] results is `390ms 500000000 rows 4.0 GB 10 GB/s`.

I've also implemented [mock server in Go](https://github.com/ClickHouse/ch-go/blob/main/internal/cmd/ch-bench-server/main.go) that simulates ClickHouse server to reduce
overhead, because currently the main bottleneck in this test is server itself (and probably localhost).
The [go-faster/ch][faster]  was able
to achieve `257ms 500000000 rows 4.0 GB 16 GB/s` which should be maximum
possible burst result, but I'm not 100% sure.

On [go-faster/ch][faster] micro-benchmarks I'm getting up to 27 GB/s, not accounting of any
network overhead (i.e. inmemory).
