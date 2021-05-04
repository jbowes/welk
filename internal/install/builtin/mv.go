package builtin

import (
	"context"
)

func Mkdir(ctx context.Context, host Host, ios IOs, args []string) error {

	return nil
}

func init() { Builtin["mkdir"] = Mkdir }
