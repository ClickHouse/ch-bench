package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-faster/errors"
	"github.com/vahid-sohrabloo/chconn"
	"github.com/vahid-sohrabloo/chconn/column"
)

func run(ctx context.Context) error {
	c, err := chconn.Connect(ctx, "")
	if err != nil {
		return errors.Wrap(err, "connect")
	}

	s, err := c.Select(ctx, "SELECT number FROM system.numbers_mt LIMIT 500000000")
	if err != nil {
		return errors.Wrap(err, "select")
	}

	start := time.Now()

	var data []uint64
	colRead := column.NewUint64(false)
	for s.Next() {
		if err := s.NextColumn(colRead); err != nil {
			return errors.Wrap(err, "column")
		}
		data = data[:0]
		colRead.ReadAllUnsafe(&data)
	}
	if err := s.Err(); err != nil {
		return errors.Wrap(err, "next")
	}

	fmt.Println(time.Since(start).Round(time.Millisecond))

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(2)
	}
}
