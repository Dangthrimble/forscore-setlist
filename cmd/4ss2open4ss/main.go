package main

import (
	"fmt"
	"os"

	"github.com/Dangthrimble/forscore-setlist/internal/cli"
	"github.com/Dangthrimble/forscore-setlist/internal/convert"
)

const name = "4ss2open4ss"

func main() {
	args := os.Args[1:]
	if len(args) == 1 && args[0] == "--version" {
		cli.Version(name)
		return
	}

	input, code, ok := cli.SingleFile(args, name, "file.4ss")
	if !ok {
		os.Exit(code)
	}

	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: read: %v\n", name, err)
		os.Exit(1)
	}
	xml, err := convert.FromNative(data, convert.SetlistTitle(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: convert: %v\n", name, err)
		os.Exit(1)
	}
	out := convert.NativeOutputPath(input)
	if err := os.WriteFile(out, xml, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "%s: write: %v\n", name, err)
		os.Exit(1)
	}
	fmt.Println(out)
}
