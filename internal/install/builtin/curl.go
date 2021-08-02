package builtin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/jbowes/sumdog/internal/install/devnull"
	"github.com/spf13/pflag"
)

func Curl(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.BoolP("silent", "s", true, "")    // ignore this one
	fs.BoolP("verbos", "v", false, "")   // ignore this one, too
	fs.BoolP("location", "L", false, "") // Ignored, as go follows Location headers, but maybe it shoudln't be

	rname := fs.BoolP("remote-name", "O", false, "")
	lname := fs.StringP("output", "o", "", "")

	outfmt := fs.StringP("write-out", "w", "", "")

	fs.BoolP("fail", "f", false, "")        // TODO: use this to make non-success responses return errors
	fs.StringSliceP("header", "H", nil, "") // TODO: use this

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

	// If we have an output string, by default drop the output.
	// the remote name flag will overwrite this.
	if *outfmt != "" {
		wc = devnull.New()
	}

	if *rname && *lname == "" {
		u, err := url.Parse(fs.Arg(0))
		if err != nil {
			return err
		}

		wc = host.Write(path.Base(u.Path))
	}
	if *lname != "" {
		wc = host.Write(*lname)
	}

	if _, err = io.Copy(wc, resp.Body); err != nil {
		return fmt.Errorf("could not copy output: %w", err)
	}

	if *outfmt == "" {
		return nil
	}

	// TODO: more robust handling here.
	out := *outfmt
	out = strings.ReplaceAll(out, "%(http_code)", fmt.Sprintf("%d", resp.StatusCode))
	_, err = ios.Out.Write([]byte(out))
	return err
}

func init() { Builtin["curl"] = Curl }
