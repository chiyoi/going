package repl

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Welcome() {
	fmt.Println("Go playground.")
}

func Usage() {
	Welcome()
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -m - Multiple statements mode.")
}

func MainLoop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	cancel := new(func() error)

	var fs flag.FlagSet
	fs.Usage = Usage

	m := fs.Bool("m", false, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return
	}

	go func() {
		Welcome()
		if *m {
			fmt.Println("Multiple statements mode.")
		} else {
			fmt.Println("Single expression mode.")
		}

		fmt.Println("Exit with ctrl-c.")
		for {
			expr, err := r(*m)
			switch err {
			case ErrExit:
				c <- syscall.SIGINT
			case ErrSkip:
				continue
			case nil:
			default:
				panic(err)
			}

			var out chan string
			out, cancel = e(expr)
			p(out)
		}
	}()

	for {
		sig := <-c
		if *cancel != nil {
			if err := (*cancel)(); err != nil {
				fmt.Println("(Error while canceling.)", err)
			}
		} else {
			fmt.Println(sig)
			break
		}
	}
	fmt.Println("Bye.")
}
