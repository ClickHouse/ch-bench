#!/bin/bash


hyperfine -w 10 -r 100 \
  ./bin/ch-bench-faster -n go-faster \
  ./bin/ch-bench-cpp -n clickhouse-cpp \
  'clickhouse-client -q "SELECT number FROM system.numbers_mt LIMIT 500000000" --format Null --time' -n clickhouse-client \
  --export-markdown RESULTS.md
