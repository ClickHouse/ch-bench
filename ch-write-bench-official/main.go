package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go"
)

func run(ctx context.Context) error {
	connect, err := clickhouse.OpenDirect("tcp://127.0.0.1:9000?username=&debug=true")
	if err != nil {
		log.Fatal(err)
	}
	{
		connect.Begin()
		stmt, _ := connect.Prepare(`
		CREATE TABLE IF NOT EXISTS test_table (id UInt8, text FixedString(4)) ENGINE = Memory
		`)

		if _, err := stmt.Exec([]driver.Value{}); err != nil {
			log.Fatal(err)
		}

		if err := connect.Commit(); err != nil {
			log.Fatal(err)
		}
	}
	v := []byte("test")
	start := time.Now()
	{
		connect.Begin()
		connect.Prepare("INSERT INTO test_table VALUES ()")

		block, err := connect.Block()
		if err != nil {
			log.Fatal(err)
		}

		block.Reserve()
		block.NumRows += 1_000_000

		for i := 0; i < 1_000_000; i++ {
			block.WriteUInt8(0, 1)
		}

		for i := 0; i < 1_000_000; i++ {
			block.WriteFixedString(1, v)
		}

		if err := connect.WriteBlock(block); err != nil {
			log.Fatal(err)
		}

		if err := connect.Commit(); err != nil {
			log.Fatal(err)
		}
	}
	duration := time.Since(start)
	fmt.Println(duration.Round(time.Millisecond))
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(2)
	}
}
