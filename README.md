# Benchmarks

Totally **unscientific** and mostly unrealistic benchmark that
[go-faster/ch](https://github.com/go-faster/ch) project uses to understand performance.

The main goal is to measure minimal **client overhead** (CPU, RAM) to read data,
i.e. data blocks deserialization and transfer.

Please see [Notes](#Notes) for more details about results.

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

## Notes

### C++
Mean results are identical and C++ has much lower dispersion:
```console
$ hyperfine ch-bench-faster ./ch-bench/ch-bench
Benchmark 1: ch-bench-faster
  Time (mean ± σ):     611.7 ms ± 133.4 ms    [User: 139.4 ms, System: 337.8 ms]
  Range (min … max):   422.2 ms … 853.1 ms    10 runs

Benchmark 2: ./ch-bench/ch-bench
  Time (mean ± σ):     676.4 ms ±  50.8 ms    [User: 445.8 ms, System: 219.6 ms]
  Range (min … max):   614.6 ms … 769.4 ms    10 runs

Summary
  'ch-bench-faster' ran
    1.11 ± 0.26 times faster than './ch-bench/ch-bench'
```

We are selecting best result, so picking `422.2 ms` vs `614.5ms`, while mean results
are much closer.

### Rust

Benchmarks were performed on `Ryzen 9 5950x`, where Rust behaves surprisingly bad:
```console
$ hyperfine ch-bench-faster ch-bench-rust
Benchmark 1: ch-bench-faster
  Time (mean ± σ):     668.6 ms ± 118.5 ms    [User: 132.6 ms, System: 392.2 ms]
  Range (min … max):   521.9 ms … 828.5 ms    10 runs

Benchmark 2: ch-bench-rust
  Time (mean ± σ):     29.703 s ±  1.614 s    [User: 28.666 s, System: 1.056 s]
  Range (min … max):   26.907 s … 31.897 s    10 runs

Summary
  'ch-bench-faster' ran
   44.43 ± 8.24 times faster than 'ch-bench-rust'
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
