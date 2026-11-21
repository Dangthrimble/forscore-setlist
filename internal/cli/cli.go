// Package cli holds the shared command-line behaviour for the setlist tools.
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/Dangthrimble/forscore-setlist/internal/version"
)

// Version prints the standard one-line version response.
func Version(name string) {
	fmt.Printf("%s %s\n", name, version.String())
}

// SingleFile parses arguments for a command that accepts exactly one input path.
// It returns the path and true when the command should continue. When it returns
// false, the caller should exit with the returned code.
func SingleFile(args []string, name, placeholder string) (string, int, bool) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Printf("usage: %s <%s>\n", name, placeholder)
		return "", 0, false
	}

	endFlags := false
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
		endFlags = true
	}
	if len(args) != 1 || (!endFlags && strings.HasPrefix(args[0], "-")) {
		fmt.Fprintf(os.Stderr, "usage: %s <%s>\n", name, placeholder)
		return "", 2, false
	}
	return args[0], 0, true
}
