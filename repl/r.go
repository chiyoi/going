package repl

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrExit = errors.New("exit")
	ErrSkip = errors.New("skip")
)

func r(multiple bool) (expr string, err error) {
	if multiple {
		expr = readLines()
	} else {
		expr = readLine()
	}
	expr = strings.Trim(expr, "\n ;")

	if expr == "e" {
		return "", ErrExit
	}
	if expr == "" {
		return "", ErrSkip
	}

	if !multiple {
		expr = fmt.Sprintf("fmt.Println(%s)", expr)
	}
	return
}

func readLines() string {
	var buf strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	fmt.Println(">--")
	for sc.Scan() && sc.Text() != "" {
		fmt.Fprintln(&buf, sc.Text())
	}
	fmt.Println("---")
	return buf.String()
}

func readLine() string {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	sc.Scan()
	return sc.Text()
}
