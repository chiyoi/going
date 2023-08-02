package repl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const tmplMain = `
package main

func main() {
%s
}
`

func e(expr string) chan string {
	out := make(chan string, 1)
	go func() {
		f, err := os.CreateTemp(os.TempDir(), "going-*.go")
		if err != nil {
			fmt.Println("Unknown error:", err)
			os.Exit(1)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		fmt.Fprintf(f, tmplMain, expr)
		f.Seek(0, io.SeekStart)

		cmd := exec.Command("gopls", "imports", f.Name())
		cmd.Stdout = f
		cmd.Stderr = ChanWriter(out)
		if err := cmd.Run(); err != nil {
			out <- fmt.Sprintf("Error occurred while preprocessing: %s\n", err)
			close(out)
			return
		}

		cmd = exec.Command("go", "run", f.Name())
		cmd.Stdout = ChanWriter(out)
		cmd.Stderr = ChanWriter(out)
		if err := cmd.Run(); err != nil {
			out <- fmt.Sprintf("Error occurred: %s\n", err)
			close(out)
			return
		}
		close(out)
	}()
	return out
}

type ChanWriter chan string

func (cw ChanWriter) Write(p []byte) (n int, err error) {
	select {
	case cw <- string(p):
		return len(p), nil
	default:
		return 0, errors.New("channel unavailable")
	}
}
