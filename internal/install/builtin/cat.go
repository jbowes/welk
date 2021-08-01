package builtin

import (
	"context"
	"io/ioutil"
)

func Cat(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: support all the args and files and stuff, instead of just stdin.

	b, err := ioutil.ReadAll(ios.In)
	if err != nil {
		return err
	}

	// TODO: send to stdout if appropriate.
	host.Log("cat", string(b))

	return nil
}

func init() { Builtin["cat"] = Cat }
