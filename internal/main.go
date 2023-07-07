package internal

import (
	"fmt"
	"os"
	"os/signal"
)

func Main() {
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
