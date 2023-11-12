package repl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const tmplMain = `
package main

func main() {
	fmt.Println(%s)
}
`

func E(ctx context.Context, expr string) (commandOut chan []byte, commandErr chan error, clean func(), err error) {
	source, err := temporarySource(expr)
	if err != nil {
		return
	}
	clean = func() {
		os.Remove(source)
	}
	defer func() {
		if err != nil {
			clean()
		}
	}()

	var errorBuffer bytes.Buffer
	cmd := exec.CommandContext(ctx, "gopls", "imports", "-w", source)
	cmd.Stderr = &errorBuffer
	err = cmd.Run()
	select {
	case <-ctx.Done():
		return
	default:
	}
	if err != nil {
		err = fmt.Errorf("preprocess error: %w\n%s", err, errorBuffer.String())
		return
	}

	errorBuffer.Reset()
	commandOut = make(chan []byte)
	commandErr = make(chan error)
	cmd = exec.CommandContext(ctx, "go", "run", source)
	cmd.Stdout = ChanWriter(commandOut)
	cmd.Stderr = &errorBuffer
	go func() {
		defer close(commandOut)
		defer close(commandErr)
		err := cmd.Run()
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err != nil {
			commandErr <- fmt.Errorf("run error: %s\n%s", err, errorBuffer.String())
			return
		}
	}()
	return
}

func temporarySource(expr string) (source string, err error) {
	f, err := os.CreateTemp("", "going-*.go")
	if err != nil {
		return
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, tmplMain, expr)
	if err != nil {
		return
	}
	return f.Name(), nil
}

var _ io.Writer = (ChanWriter)(nil)

type ChanWriter chan []byte

func (c ChanWriter) Write(p []byte) (n int, err error) {
	c <- p
	return len(p), nil
}
