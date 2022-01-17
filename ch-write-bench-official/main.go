package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/dustin/go-humanize"
)

func run(ctx context.Context) error {
	c, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return err
	}
	if err := c.Exec(ctx, "CREATE TABLE IF NOT EXISTS test_table (id UInt64) ENGINE = Null"); err != nil {
		return err
	}
	start := time.Now()
	const (
		totalBlocks = 500
		rowsInBlock = 1_000_000
		totalRows   = totalBlocks * rowsInBlock
		totalBytes  = totalRows * (64 / 8)
	)
	var (
		idColumns []uint64
	)
	for i := 0; i < rowsInBlock; i++ {
		idColumns = append(idColumns, 1)
	}
	{
		for i := 0; i < totalBlocks; i++ {
			batch, err := c.PrepareBatch(ctx, "INSERT INTO test_table VALUES")
			if err != nil {
				return err
			}
			if err := batch.Column(0).Append(idColumns); err != nil {
				return err
			}
			if err := batch.Send(); err != nil {
				return err
			}
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
