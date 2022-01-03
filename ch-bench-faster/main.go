package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-faster/ch"
	"github.com/go-faster/ch/proto"
	"github.com/go-faster/errors"
)

func run(ctx context.Context) error {
	c, err := ch.Dial(ctx, "localhost:9000", ch.Options{})
	if err != nil {
		return errors.Wrap(err, "dial")
	}
	defer func() { _ = c.Close() }()

	var (
		gotRows  uint64
		gotBytes uint64
		data     proto.ColUInt64
	)
	start := time.Now()
	if err := c.Do(ctx, ch.Query{
		Body: "SELECT number FROM system.numbers_mt LIMIT 500000000",
		OnProgress: func(ctx context.Context, p proto.Progress) error {
			gotBytes += p.Bytes
			return nil
		},
		OnResult: func(ctx context.Context, block proto.Block) error {
			gotRows += uint64(block.Rows)
			return nil
		},
		Result: proto.Results{
			{Name: "number", Data: &data},
		},
	}); err != nil {
		return errors.Wrap(err, "query")
	}

	duration := time.Since(start)
	fmt.Println(duration.Round(time.Millisecond), gotRows, "rows",
		humanize.Bytes(gotBytes),
		humanize.Bytes(uint64(float64(gotBytes)/duration.Seconds()))+"/s",
	)

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(2)
	}
}
