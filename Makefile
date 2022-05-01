.PHONY: ch-bench-chconn
.PHONY: ch-bench-uptrace
.PHONY: ch-bench-faster
.PHONY: ch-bench-rust
.PHONY: ch-bench-rust-driver
.PHONY: ch-bench-rust-http
.PHONY: ch-bench-official
.PHONY: ch-bench-mailru
.PHONY: build

ch-bench-chconn:
	go build -o bin ./ch-bench-chconn

ch-bench-faster:
	go build -o bin ./ch-bench-faster

ch-bench-uptrace:
	go build -o bin ./ch-bench-uptrace

ch-bench-official:
	go build -o bin ./ch-bench-official

ch-bench-mailru:
	go build -o bin ./ch-bench-mailru

ch-bench-rust:
	cd ch-bench-rust && RUSTFLAGS="-C target-cpu=native" cargo build --profile release-adjusted
	rm -f ./bin/ch-bench-rust
	cp ./ch-bench-rust/target/release/ch-bench-rust ./bin/ch-bench-rust

ch-bench-rust-driver:
	cd ch-bench-rust-driver && RUSTFLAGS="-C target-cpu=native" cargo build --profile release-adjusted
	rm -f ./bin/ch-bench-rust-driver
	cp ./ch-bench-rust-driver/target/release/ch-bench-rust-driver ./bin/ch-bench-rust-driver

ch-bench-rust-http:
	cd ch-bench-rust-http && RUSTFLAGS="-C target-cpu=native" cargo build --profile release-adjusted
	rm -f ./bin/ch-bench-rust-http
	cp ./ch-bench-rust-http/target/release/ch-bench-rust-http ./bin/ch-bench-rust-http

build: ch-bench-chconn ch-bench-official ch-bench-faster ch-bench-uptrace ch-bench-rust ch-bench-rust-http ch-bench-mailru ch-bench-official ch-bench-rust-driver

run:
	hyperfine -w 10 -r 100 \
	  ./bin/ch-bench-faster -n go-faster \
	  ./bin/ch-bench-cpp -n clickhouse-cpp \
	  ./bin/ch-bench-chconn -n vahid-sohrabloo/chconn \
	  ./bin/ch-bench-rust-driver -n clickhouse_driver_rust \
	  ./bin/ch-bench-official -n clickhouse-go  \
	  ./bin/ch-bench-uptrace -n uptrace  \
	  'clickhouse-client -q "SELECT number FROM system.numbers_mt LIMIT 500000000" --format Null --time' -n clickhouse-client \
	  --export-markdown RESULTS.md
run-slow:
	hyperfine -r 5 \
	  ./bin/ch-bench-faster -n go-faster \
	  ./bin/ch-bench-cpp -n cpp \
	  ./bin/ch-bench-rust -n rs \
	  ./bin/ch-bench-rust-http -n rs-http \
	  ./bin/ch-bench-chconn -n vahid-sohrabloo/chconn \
	  'clickhouse-client -q "SELECT number FROM system.numbers_mt LIMIT 500000000" --format Null --time' -n clickhouse-client \
	  --export-markdown RESULTS.slow.md
