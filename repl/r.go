package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chiyoi/apricot/logs"
)

const (
	ps1 = ">> "
	ps2 = "-> "
)

func R(ctx context.Context, input chan []byte, showPrompt bool) (expr string, err error) {
	push := func(stack *[]byte, b byte) {
		*stack = append(*stack, b)
	}

	pop := func(stack *[]byte) (b byte) {
		if len(*stack) == 0 {
			return
		}
		*stack, b = (*stack)[:len(*stack)-1], (*stack)[len(*stack)-1]
		return
	}

	pairBracket := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	var lineBuffer strings.Builder
	if showPrompt {
		fmt.Fprint(os.Stdout, ps1)
	}

	var bracketStack []byte
	for {
		var bs []byte
		select {
		case <-ctx.Done():
			return
		case bs = <-input:
		}

		if bs == nil {
			if lineBuffer.Len() == 0 {
				err = ErrExit
				return
			}
			break
		}

		var flag bool
		if len(bs) >= 3 {
			switch string(bs[:3]) {
			case ps1, ps2:
				flag = true
				bs = bs[3:]
			}
		}

		lineBuffer.Write(bs)
		lineBuffer.WriteByte('\n')
		for _, b := range bs {
			switch b {
			case '(', '{', '[':
				push(&bracketStack, b)
			case ')', '}', ']':
				if pairBracket[b] != pop(&bracketStack) {
					err = fmt.Errorf("%w: parsing `%c`", ErrInvalidSyntax, b)
					return
				}
			}
		}

		if len(bracketStack) == 0 {
			break
		}
		if !flag && showPrompt {
			fmt.Fprint(os.Stdout, ps2)
			fmt.Fprint(os.Stdout, strings.Repeat("\t", len(bracketStack)))
		}
	}

	line := lineBuffer.String()
	line = strings.Trim(line, " \n\t;")
	return line, nil
}

func StartScanner() (input chan []byte) {
	input = make(chan []byte)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for {
			if !sc.Scan() {
				input <- nil
				sc = bufio.NewScanner(os.Stdin)
				continue
			}
			logs.Debug("Scanned.", "sc.Bytes():", sc.Bytes())
			input <- sc.Bytes()
		}
	}()
	return input
}
