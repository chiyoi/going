package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func r() string {
	var buf strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	fmt.Println("<--")
	for sc.Scan() {
		fmt.Fprintln(&buf, sc.Text())
	}
	fmt.Println("-->")
	return buf.String()
}
