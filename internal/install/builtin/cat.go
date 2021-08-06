package builtin

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/jbowes/welk/internal/install/devnull"
)

func Cat(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: support all the args and files and stuff, instead of just stdin.

	var b []byte
	var err error

	if len(args) == 0 {
		b, err = ioutil.ReadAll(ios.In)
		if err != nil {
			return err
		}
	} else {
		for _, f := range args {
			// used for a multi-line comment trick in goreleaser / go downloader scripts
			if f == "/dev/null" {
				continue
			}
			fmt.Println(f)
			// TODO: check for existance and error?
			b = append(b, host.File(ctx, f)...)
		}
	}

	// TODO: this isn't accurate. it should log when going to the controlling
	// console (which happens do be IsDevNull right now).
	if devnull.IsDevNull(ios.Out) {
		host.Log("cat", string(b))
	} else {
		io.WriteString(ios.Out, string(b))
	}

	return nil
}

func init() { Builtin["cat"] = Cat }
