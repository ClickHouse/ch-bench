package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func run(ctx context.Context) error {
	c, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return err
	}

	start := time.Now()
	rows, err := c.Query(ctx, "SELECT number FROM system.numbers_mt LIMIT 500000000")
	if err != nil {
		return err
	}
	var count int
	for rows.Next() {
		count++
	}

	fmt.Println(time.Since(start).Round(time.Millisecond), count)

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(2)
	}
}
