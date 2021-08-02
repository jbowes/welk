package builtin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"

	"github.com/spf13/pflag"
)

func Tar(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: supported bundled args? eg tar zxf instead of tar -zxf

	fs := pflag.NewFlagSet("", pflag.ContinueOnError)

	fs.Bool("no-same-owner", false, "") // ignored

	gz := fs.BoolP("z", "z", false, "")
	extract := fs.BoolP("x", "x", false, "")
	fname := fs.StringP("f", "f", "", "")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	host.Log("tar")

	if !*extract {
		return errors.New("tar: only extract is supported")
	}

	in := ios.In
	if *fname != "" {
		in = bytes.NewReader(host.File(*fname))
	}

	if *gz {
		in, err = gzip.NewReader(in)
		if err != nil {
			return err
		}
	}

	tr := tar.NewReader(in)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		if hdr.Typeflag == tar.TypeDir {
			host.MkDir(hdr.Name)
			continue
		}

		err = func() error {
			o := host.Write(hdr.Name)
			defer o.Close()

			_, err := io.Copy(o, tr)
			return err
		}() // get that defer
		if err != nil {
			return err
		}
	}

	return nil
}

func init() { Builtin["tar"] = Tar }
