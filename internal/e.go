package internal

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
%s
}
`

func e(expr string) (context.Context, io.Reader) {
	ctx, cancel := context.WithCancel(context.Background())
	var buf bytes.Buffer

	go func() {
		f, err := os.CreateTemp(os.TempDir(), "going-*.go")
		if err != nil {
			fmt.Println("Unknown error:", err)
			os.Exit(1)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		fmt.Fprintf(
			f,
			tmplMain,
			expr,
		)
		f.Seek(0, io.SeekStart)

		cmd := exec.Command("gopls", "imports", f.Name())
		cmd.Stdout = f
		cmd.Stderr = &buf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(
				&buf,
				"Error occurred while preprocessing: %s\nOutput: %s\n",
				err,
				&buf,
			)
			cancel()
			return
		}

		cmd = exec.Command("go", "run", f.Name())
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(
				&buf,
				"Error occurred: %s\nOutput:\n%s",
				err,
				&buf,
			)
		}
		cancel()
	}()
	return ctx, &buf
}
