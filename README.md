# Benchmarks

Totally **unscientific** and mostly unrealistic benchmark that
[go-faster/ch](https://github.com/go-faster/ch) project uses to understand performance.

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

| Name                                          | Time    | RAM    | Speedup |
|-----------------------------------------------|---------|--------|---------|
| clickhouse-client (C++)                       | 393ms   | 91M    | ~1x     |
| **go-faster/ch**                              | 395ms   | 9M     | 1x      |
| clickhouse-cpp (C++)                          | 531ms   | 6.9M   | 1.34x   |
| *clickhouse-rs (Rust, AMD EPYC **Adjusted**)* | *740ms* | *192M* | *1.68x* |
| vahid-sohrabloo/chconn (Go)                   | 5s      | 10M    | 11x     |
| clickhouse-jdbc (Java)                        | 10s     | 702M   | 22x     |
| clickhouse-rs (Rust, AMD Ryzen 9)             | 27s     | 192M   | 61x     |
| clickhouse-go                                 | 35s     | 184M   | 79x     |
| clickhouse-driver (Python)                    | 37s     | 60M    | 84x     |
| mailru/go-clickhouse                          | 4m13s   | 13M    | 575x    |

NB: **mailru/go-clickhouse** and **clickhouse-jdbc** are using HTTP protocol.

## Notes

### C++
Mean results are identical and C++ has much lower dispersion:

| Command             |     Mean [ms] | Min [ms] | Max [ms] |    Relative |
|:--------------------|--------------:|---------:|---------:|------------:|
| `clickhouse-cpp`    |  575.2 ± 36.5 |    531.3 |    686.1 |        1.00 |
| `clickhouse-client` | 611.5 ± 161.1 |    393.2 |   1102.6 | 1.06 ± 0.29 |
| `go-faster`         |  626.4 ± 90.9 |    395.5 |    805.1 | 1.09 ± 0.17 |


We are selecting **best** result, so picking `393 ms` vs `531 ms`, while mean results
are much closer.

### Rust

Benchmarks were performed on `Ryzen 9 5950x`, where Rust behaves surprisingly bad:
```console
Benchmark 1: go-faster
  Time (mean ± σ):     644.6 ms ±  53.8 ms    [User: 109.7 ms, System: 352.5 ms]
  Range (min … max):   586.8 ms … 719.4 ms    5 runs

Benchmark 2: clickhouse-cpp
  Time (mean ± σ):     579.5 ms ±  23.2 ms    [User: 381.7 ms, System: 185.1 ms]
  Range (min … max):   541.8 ms … 599.0 ms    5 runs

Benchmark 3: clickhouse-rs
  Time (mean ± σ):     27.122 s ±  1.342 s    [User: 26.129 s, System: 1.024 s]
  Range (min … max):   24.760 s … 28.106 s    5 runs

  Warning: Statistical outliers were detected. Consider re-running this benchmark on a quiet PC without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options.

Benchmark 4: vahid-sohrabloo/chconn
  Time (mean ± σ):      5.066 s ±  0.115 s    [User: 4.632 s, System: 0.535 s]
  Range (min … max):    4.901 s …  5.204 s    5 runs

Benchmark 5: clickhouse-go
  Time (mean ± σ):     38.254 s ±  0.098 s    [User: 74.100 s, System: 1.179 s]
  Range (min … max):   38.120 s … 38.366 s    5 runs

Benchmark 6: clickhouse-client
  Time (mean ± σ):     507.6 ms ±  97.2 ms    [User: 135.3 ms, System: 197.7 ms]
  Range (min … max):   408.5 ms … 615.7 ms    5 runs
```

However, on Intel results are much closer:
```console
Benchmark 1: ch-bench-rust
  Time (mean ± σ):      5.309 s ±  1.845 s    [User: 4.852 s, System: 0.727 s]
  Range (min … max):    2.055 s …  8.683 s    10 runs

Benchmark 2: ch-bench-faster
  Time (mean ± σ):      1.435 s ±  0.138 s    [User: 0.364 s, System: 0.767 s]
  Range (min … max):    1.122 s …  1.588 s    10 runs

Summary
  'ch-bench-faster' ran
    3.70 ± 1.33 times faster than 'ch-bench-rust'
```

Also, on AMD EPYC they are even closer:
```console
$ hyperfine ch-bench-rust ch-bench-faster
Benchmark 1: ch-bench-rust
  Time (mean ± σ):      3.949 s ±  1.324 s    [User: 2.133 s, System: 2.188 s]
  Range (min … max):    2.672 s …  6.198 s    10 runs

Benchmark 2: ch-bench-faster
  Time (mean ± σ):      2.020 s ±  0.091 s    [User: 0.348 s, System: 1.399 s]
  Range (min … max):    1.893 s …  2.225 s    10 runs

Summary
  'ch-bench-faster' ran
    1.95 ± 0.66 times faster than 'ch-bench-rust'
```

Please create an issue to help me improve results on `Ryzen 9 5950x` if it is possible,
Rust client is pretty good and should perform better.
