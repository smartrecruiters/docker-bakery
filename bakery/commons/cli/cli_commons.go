// Package cli is for common CLI utilities
package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

// CustomExitHandler checks if the error fulfills the ExitCoder interface, and if
// so prints the error to stderr (if it is non-empty) and calls OsExiter with the
// given exit code.  If the given error is a MultiError, then this func is
// called on all members of the Errors slice and calls OsExiter with the last exit code.
func CustomExitHandler(context *cli.Context, err error) {
	// invoke default handler first
	cli.HandleExitCoder(err)
	// SR: This is the actual change done in HandleExitCoder to print error even if it is not the MultiError or ExitCoder
	if err != nil && err.Error() != "" {
		fmt.Fprintln(cli.ErrWriter, err)
		cli.OsExiter(1)
	}
}
