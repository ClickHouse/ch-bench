package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

func run(ctx context.Context) error {
	c, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000")
	if err != nil {
		return err
	}

	start := time.Now()
	rows, err := c.QueryContext(ctx, "SELECT number FROM system.numbers LIMIT 500000000")
	if err != nil {
		return err
	}
	var count int
	for rows.Next() {
		var v uint64
		if err := rows.Scan(&v); err != nil {
			return err
		}
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
