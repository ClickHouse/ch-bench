package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	createTable := ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS test_table (id UInt8, text FixedString(4)) ENGINE = Memory",
	}
	if err := c.Do(ctx, createTable); err != nil {
		return err
	}
	start := time.Now()
	var (
		idColumns proto.ColUInt8

		v           = []byte("test")
		textColumns = proto.ColFixedStr{Size: 4}
	)
	for i := 0; i < 1_000_000; i++ {
		idColumns = append(idColumns, 1)
		textColumns.Buf = append(textColumns.Buf, v...)
	}
	insertQuery := ch.Query{
		Body: "INSERT INTO test_table VALUES",
		Input: []proto.InputColumn{
			{Name: "id", Data: &idColumns},
			{Name: "text", Data: &textColumns},
		},
	}

	if err := c.Do(ctx, insertQuery); err != nil {
		return err
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
