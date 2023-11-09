package command

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	OperationCancelledString = "command cancelled by user"
)

// ConfirmAction now takes in input and output streams.
func (ctx *Context) ConfirmAction(prompt string, in io.Reader, out io.Writer) (bool, error) {

	if ctx.cfg.DefaultOverride {
		return true, nil
	}

	scanner := bufio.NewScanner(in)
	for {
		_, err := fmt.Fprint(out, prompt+" [y/n]: ")
		if err != nil {
			return false, fmt.Errorf("error writing to output: %w", err)
		}
		if !scanner.Scan() {
			if scanner.Err() != nil {
				return false, fmt.Errorf("error reading input: %w", scanner.Err())
			}
			return false, fmt.Errorf(OperationCancelledString)
		}

		response := scanner.Text()
		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			_, err := fmt.Fprintln(out, "Invalid input. Please enter 'y' for yes or 'n' for no.")
			if err != nil {
				return false, fmt.Errorf("error writing to output: %w", err)
			}
		}
	}
}
