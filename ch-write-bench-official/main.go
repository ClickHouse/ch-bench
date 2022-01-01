package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/dustin/go-humanize"
)

func run(ctx context.Context) error {
	connect, err := clickhouse.OpenDirect("tcp://127.0.0.1:9000?username=&debug=true")
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	const (
		totalBlocks = 500
		rowsInBlock = 1_000_000
		totalRows   = totalBlocks * rowsInBlock
		totalBytes  = totalRows * (64 / 8)
	)
	{
		connect.Begin()
		connect.Prepare("INSERT INTO test_table VALUES ()")

		for i := 0; i < totalBlocks; i++ {
			block, err := connect.Block()
			if err != nil {
				log.Fatal(err)
			}
			block.Reserve()
			block.NumRows += rowsInBlock
			for i := 0; i < rowsInBlock; i++ {
				block.WriteUInt64(0, 1)
			}
			if err := connect.WriteBlock(block); err != nil {
				log.Fatal(err)
			}
		}
		if err := connect.Commit(); err != nil {
			log.Fatal(err)
		}
	}
	duration := time.Since(start)
	fmt.Println(duration.Round(time.Millisecond), totalRows, "rows",
		humanize.Bytes(totalBytes),
		humanize.Bytes(uint64(float64(totalBytes)/duration.Seconds()))+"/s",
	)
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(2)
	}
}
