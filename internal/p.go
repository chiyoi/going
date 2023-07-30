package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

func p(ctx context.Context, r io.Reader) {
	buf := make([]byte, 64)
	for {
		_, err := r.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println("Unknown error:", err)
			return
		}
		log.Println("(Read returned.)", err)
		fmt.Print(string(buf))
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 100):
		}

	}
}
