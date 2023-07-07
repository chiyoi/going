package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
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
