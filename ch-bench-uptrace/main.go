package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/uptrace/go-clickhouse/ch"
)

func run(ctx context.Context) error {
	db := ch.Connect(ch.WithTimeout(time.Second * 30))

	start := time.Now()
	rows, err := db.QueryContext(ctx, "SELECT number FROM system.numbers_mt LIMIT 500000000")
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
