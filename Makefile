.PHONY: ch-bench-chconn
.PHONY: ch-bench-faster
.PHONY: ch-bench-rust
.PHONY: ch-bench-official
.PHONY: ch-bench-mailru
.PHONY: build

ch-bench-chconn:
	go build -o bin ./ch-bench-chconn

ch-bench-faster:
	go build -o bin ./ch-bench-faster

ch-bench-official:
	go build -o bin ./ch-bench-official

ch-bench-mailru:
	go build -o bin ./ch-bench-mailru

ch-bench-rust:
	cd ch-bench-rust && cargo build --release
	rm ./bin/ch-bench-rust
	cp ./ch-bench-rust/target/release/ch-bench-rust ./bin/ch-bench-rust

build: ch-bench-chconn ch-bench-faster ch-bench-rust ch-bench-mailru ch-bench-official

run:
	hyperfine -w 10 -r 100 \
	  ./bin/ch-bench-faster -n go-faster \
	  ./bin/ch-bench-cpp -n clickhouse-cpp \
	  'clickhouse-client -q "SELECT number FROM system.numbers_mt LIMIT 500000000" --format Null --time' -n clickhouse-client \
	  --export-markdown RESULTS.md
run-slow:
	hyperfine -r 5 \
	  ./bin/ch-bench-faster -n go-faster \
	  ./bin/ch-bench-cpp -n clickhouse-cpp \
	  ./bin/ch-bench-rust -n clickhouse-rs \
	  ./bin/ch-bench-chconn -n vahid-sohrabloo/chconn \
	  ./bin/ch-bench-official -n clickhouse-go  \
	  'clickhouse-client -q "SELECT number FROM system.numbers_mt LIMIT 500000000" --format Null --time' -n clickhouse-client \
	  --export-markdown RESULTS.slow.md
