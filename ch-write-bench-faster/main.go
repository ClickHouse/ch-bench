package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/dustin/go-humanize"
	"github.com/go-faster/errors"
)

func run(ctx context.Context) error {
	c, err := ch.Dial(ctx, ch.Options{
		Compression: ch.CompressionLZ4,
	})
	if err != nil {
		return errors.Wrap(err, "dial")
	}
	defer func() { _ = c.Close() }()

	if err := c.Do(ctx, ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS test_table (id UInt64) ENGINE = Null",
	}); err != nil {
		return err
	}
	start := time.Now()
	const (
		totalBlocks = 5000
		rowsInBlock = 60_000
		totalRows   = totalBlocks * rowsInBlock
		totalBytes  = totalRows * (64 / 8)
	)
	var (
		idColumns proto.ColUInt64
		blocks    int
	)
	for i := 0; i < rowsInBlock; i++ {
		idColumns = append(idColumns, 1)
	}
	if err := c.Do(ctx, ch.Query{
		Body: "INSERT INTO test_table VALUES",
		OnInput: func(ctx context.Context) error {
			blocks++
			if blocks >= totalBlocks {
				return io.EOF
			}
			return nil
		},
		Input: []proto.InputColumn{
			{Name: "id", Data: idColumns},
		},
	}); err != nil {
		return err
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
