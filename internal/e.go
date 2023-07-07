package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const tmplMain = `
package main

func main() {
%s
}
`

func e(expr string) string {
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

	var buf strings.Builder
	cmd := exec.Command("gopls", "imports", f.Name())
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf(
			"Error occurred while preprocessing: %s\nOutput: %s\n",
			err,
			&buf,
		)
	}

	if err := func() (err error) {
		if _, err = f.Seek(0, io.SeekStart); err != nil {
			return
		}
		return f.Truncate(0)
	}(); err != nil {
		fmt.Println("Unknown error:", err)
		os.Exit(1)
	}
	fmt.Fprint(f, buf.String())

	buf.Reset()
	cmd = exec.Command("go", "run", f.Name())
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf(
			"Error occurred: %s\nOutput:\n%s",
			err,
			&buf,
		)
	}

	return buf.String()
}
