package internal

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func Usage() {
	fmt.Println("Go playground.")
	fmt.Println("Code your program (end with ctrl-d), and I will")
	fmt.Println("fill it into the main function, add necessary")
	fmt.Println("imports and run it.")
}

func Main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var fs flag.FlagSet
	fs.Usage = Usage
	if err := fs.Parse(os.Args[1:]); err != nil {
		return
	}

	go func() {
		Usage()
		fmt.Println("Exit with ctrl-c.")
		for {
			expr := r()
			res := e(expr)
			p(res)
		}
	}()

	sig := <-c
	fmt.Println(sig)
	fmt.Println("Bye.")
}
