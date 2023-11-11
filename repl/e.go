package repl

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

const tmplMain = `
package main

func main() {
	fmt.Println(%s)
}
`

func E(ctx context.Context, expr string) (output string, err error) {
	source, err := temporarySource(expr)
	if err != nil {
		return
	}
	defer os.Remove(source)

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
	cmd = exec.CommandContext(ctx, "go", "run", source)
	cmd.Stderr = &errorBuffer
	bs, err := cmd.Output()
	output = string(bs)
	select {
	case <-ctx.Done():
		return
	default:
	}
	if err != nil {
		err = fmt.Errorf("run error: %w\n%s", err, errorBuffer.String())
		return
	}
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
