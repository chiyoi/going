package repl

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func Welcome() {
	fmt.Println("Go playground.")
	fmt.Println("Code your program, and I will fill it into the main function,")
	fmt.Println("add necessary imports and run it.")
}

func Usage() {
	Welcome()
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -e, -expr - Single expression mode.")
}

func MainLoop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var fs flag.FlagSet
	fs.Usage = Usage

	em := fs.Bool("expr", false, "")
	fs.BoolVar(em, "e", false, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return
	}

	go func() {
		Welcome()
		fmt.Println("Exit with ctrl-c.")
		for {
			var expr string
			if *em {
				expr = re()
			} else {
				expr = r()
			}

			res := e(expr)
			p(res)
		}
	}()

	sig := <-c
	fmt.Println(sig)
	fmt.Println("Bye.")
}
