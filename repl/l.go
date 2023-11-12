package repl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/chiyoi/apricot/logs"
)

func L(input chan []byte, showPrompt bool) {
	ctx := context.Background()
	for {
		err := work(ctx, input, showPrompt)
		logs.Debug(err)
		switch {
		case err == nil, errors.Is(err, ErrSkip):
		case errors.Is(err, ErrExit):
			fmt.Println("Bye.")
			return
		case errors.Is(err, ErrContinue):
			fmt.Println()
		default:
			fmt.Println(err)
		}
	}
}

func work(ctx context.Context, input chan []byte, showPrompt bool) (err error) {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	expr, err := R(signalCtx, input, showPrompt)
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

	commandOut, commandErr, clean, err := E(signalCtx, expr)
	select {
	case <-signalCtx.Done():
		return ErrContinue
	default:
	}
	if err != nil {
		return
	}
	defer clean()

	err = P(commandOut, commandErr)
	select {
	case <-signalCtx.Done():
		return ErrContinue
	default:
	}
	if err != nil {
		return
	}
	return
}
