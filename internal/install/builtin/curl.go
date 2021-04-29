package builtin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/spf13/pflag"
)

func Curl(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.BoolP("silent", "s", true, "")    // ignore this one
	fs.BoolP("verbos", "v", false, "")   // ignore this one, too
	fs.BoolP("location", "L", false, "") // Ignored, as go follows Location headers, but maybe it shoudln't be

	rname := fs.BoolP("remote-name", "O", false, "") // TODO: use this to output the file
	fs.BoolP("fail", "f", false, "")                 // TODO: use this to make non-success responses return errors
	fs.StringSliceP("header", "H", nil, "")          // TODO: use this
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("curl", fs.Arg(0))

	resp, err := http.Get(fs.Arg(0))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wc := ios.Out
	if *rname {
		u, err := url.Parse(fs.Arg(0))
		if err != nil {
			return err
		}

		wc = host.Write(path.Base(u.Path))
	}

	_, err = io.Copy(wc, resp.Body)
	return err
}

func init() { Builtin["curl"] = Curl }
