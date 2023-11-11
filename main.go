package main

import (
	"os"

	"github.com/chiyoi/apricot/sakana"
	"github.com/chiyoi/going/repl"
)

func main() {
	c := repl.Handler()
	c.ServeArgs(sakana.Files{}, os.Args[1:])
}
