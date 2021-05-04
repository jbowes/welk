package install

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jbowes/sumdog/internal/install/builtin"
	"github.com/jbowes/sumdog/internal/install/store"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Run(ctx context.Context, permittedExec func([]string) bool, log func(string, ...string), url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := syntax.NewParser().Parse(resp.Body, "")
	if err != nil {
		return err
	}

	// Replace the builtins with our own builtins
	syntax.Walk(f, func(n syntax.Node) bool {
		switch x := n.(type) {
		case *syntax.CallExpr:

			if len(x.Args) > 0 && x.Args[0].Lit() == "echo" {
				x.Args[0].Parts[0].(*syntax.Lit).Value = "sumdog-echo"
			}

			if len(x.Args) > 0 && x.Args[0].Lit() == "printf" {
				x.Args[0].Parts[0].(*syntax.Lit).Value = "sumdog-printf"
			}

		}
		return true
	})

	run := &runner{
		builtin:       builtin.Builtin,
		store:         &store.Store{},
		permittedExec: permittedExec,
		log:           log,
	}

	int, err := interp.New(
		// interp.Dir(/* what makes sense here? */),
		interp.Env(nil), // TODO: configurable inclusion list
		interp.ExecHandler(run.ExecHandler),
		interp.OpenHandler(run.OpenHandler),
		// interp.Params(), /* passed in by user */
		// interp.StdIO(), // capture and log output and error
	)
	if err != nil {
		return err
	}

	err = int.Run(ctx, f)
	if err != nil {
		return err
	}

	syntax.NewPrinter().Print(os.Stdout, f)

	return nil
}

type runner struct {
	builtin map[string]builtin.BuiltinFunc
	store   *store.Store

	permittedExec func(args []string) bool
	log           func(tag string, msg ...string)
}

func (r *runner) Log(tag string, msg ...string)    { r.log(tag, msg...) }
func (r *runner) ChDir(path string)                { r.store.ChDir(path) }
func (r *runner) Write(path string) io.WriteCloser { return r.store.Write(path) }
func (r *runner) MkDir(path string)                { r.store.MkDir(path) }

func (r *runner) ExecHandler(ctx context.Context, args []string) error {
	b, ok := r.builtin[args[0]]
	if !ok {
		if r.permittedExec(args) {
			return interp.DefaultExecHandler(2)(ctx, args)
		}

		return fmt.Errorf("unimplemented command: %s", args[0])
	}

	hc := interp.HandlerCtx(ctx)
	return b(ctx, r, builtin.IOs{Out: hc.Stdout}, args[1:])
}

func (r *runner) OpenHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return devNull{}, nil
	}

	return nil, fmt.Errorf("shell file opening not implemented")
}
