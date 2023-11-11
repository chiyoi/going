package repl

import (
	"context"
	"os"
	"os/signal"
)

func L(ctx context.Context, input chan []byte) (err error) {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	expr, err := R(signalCtx, input, true)
	select {
	case <-signalCtx.Done():
		return ErrContinue
	default:
	}
	if err != nil {
		return
	}

	switch expr {
	case "e", "exit":
		return ErrExit
	case "":
		return ErrSkip
	}

	output, err := E(signalCtx, expr)
	select {
	case <-signalCtx.Done():
		return ErrContinue
	default:
	}
	if err != nil {
		return
	}

	P(output)
	return
}
