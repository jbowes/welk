package builtin

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
)

func Openssl(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: assumes dgst -sha 256
	if len(args) != 3 || args[0] != "dgst" || args[1] != "-sha256" {
		return errors.New("only openssl dgst -sha256 is supported")
	}

	s := sha256.Sum256(host.File(ctx, args[2]))
	_, err := fmt.Fprintf(ios.Out, "SHA256(%s)= %x\n", args[2], s)

	return err
}

func init() { Builtin["openssl"] = Openssl }
