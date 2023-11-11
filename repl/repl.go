package repl

import (
	"context"
	"errors"
	"fmt"

	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/sakana"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
	ErrExit          = errors.New("exit")
	ErrContinue      = errors.New("continue")
	ErrSkip          = errors.New("skip")
)

func Handler() sakana.Handler {
	input := StartScanner()
	c := sakana.NewCommand("going", "going", "Go playground.")
	c.Work(sakana.HandlerFunc(func(f sakana.Files, args []string) int {
		ctx := context.Background()
		for {
			err := L(ctx, input)
			logs.Debug(err)
			switch {
			case err == nil, errors.Is(err, ErrSkip):
			case errors.Is(err, ErrExit):
				fmt.Println("Bye.")
				return 0
			case errors.Is(err, ErrContinue):
				fmt.Println()
			default:
				fmt.Println(err)
			}
		}
	}))
	return c
}
