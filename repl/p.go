package repl

import (
	"fmt"
)

func P(commandOut chan []byte, commandErr chan error) error {
	for {
		select {
		case out, more := <-commandOut:
			if !more {
				return nil
			}
			fmt.Print(string(out))
		case err := <-commandErr:
			return err
		}
	}
}
