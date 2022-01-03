.PHONY: ch-bench-chconn
.PHONY: ch-bench-faster
.PHONY: ch-bench-rust
.PHONY: build

ch-bench-chconn:
	go build -o bin ./ch-bench-chconn

ch-bench-faster:
	go build -o bin ./ch-bench-faster

ch-bench-rust:
	cd ch-bench-rust && cargo build --release
	mv ./ch-bench-rust/target/release/ch-bench-rust ./bin/

build: ch-bench-chconn ch-bench-faster ch-bench-rust

run:
	./run.sh
