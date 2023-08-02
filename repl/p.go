package repl

import (
	"fmt"
)

func p(res chan string) {
	for s := range res {
		fmt.Print(s)
	}
}
