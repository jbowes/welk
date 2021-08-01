package builtin

import (
	"context"
)

func Ln(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("ln")

	// TODO: actually implement this
	return nil
}

func init() { Builtin["ln"] = Ln }
