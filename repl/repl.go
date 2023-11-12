package repl

import (
	"errors"
	"fmt"
	"os"

	"github.com/chiyoi/apricot/sakana"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
	ErrExit          = errors.New("exit")
	ErrContinue      = errors.New("continue")
	ErrSkip          = errors.New("skip")
)

func Handler() sakana.Handler {
	c := sakana.NewCommand("going", "going", "Go playground.")
	c.OptionUsage([]string{"h", "help"}, false, "Show this help message.")

	script := c.FlagSet.String("s", "", "")
	c.FlagSet.StringVar(script, "script", "", "")
	c.OptionUsage([]string{"s", "script"}, false, "Run a script.")
	c.Work(sakana.HandlerFunc(func(f sakana.Files, args []string) int {
		if *script != "" {
			f, err := os.Open(*script)
			if err != nil {
				fmt.Println(err)
				return 1
			}
			input := Scan(f)
			L(input, false)
			return 0
		}
		return sakana.Continue
	}))

	c.Work(sakana.HandlerFunc(func(f sakana.Files, args []string) int {
		input := Scan(os.Stdin)
		L(input, true)
		return 0
	}))
	return c
}
