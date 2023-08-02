package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func r() string {
	var buf strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	fmt.Println(">--")
	for sc.Scan() && sc.Text() != "" {
		fmt.Fprintln(&buf, sc.Text())
	}
	fmt.Println("---")
	return buf.String()
}

func re() string {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	sc.Scan()
	return fmt.Sprintf("fmt.Println(%s)", sc.Text())
}
