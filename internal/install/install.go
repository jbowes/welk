package install

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/jbowes/sumdog/internal/db"
	"github.com/jbowes/sumdog/internal/install/builtin"
	"github.com/jbowes/sumdog/internal/install/devnull"
	"github.com/jbowes/sumdog/internal/install/sham"
	"github.com/jbowes/sumdog/internal/install/vfs"
	"mvdan.cc/sh/v3/expand"
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

			// TODO: add printf impl
			if len(x.Args) > 0 && x.Args[0].Lit() == "printf" {
				x.Args[0].Parts[0].(*syntax.Lit).Value = "sumdog-printf"
			}

			if len(x.Args) > 0 && x.Args[0].Lit() == "cd" {
				x.Args[0].Parts[0].(*syntax.Lit).Value = "sumdog-cd"
			}

		}
		return true
	})

	v := &vfs.VFS{}
	run := &runner{
		builtin:       builtin.Builtin,
		sham:          sham.Sham,
		permittedExec: permittedExec,
		vfs:           v,
		log:           log,
	}

	homevarU, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	homevar := homevarU.String()

	dn := devnull.New()

	int, err := interp.New(
		// interp.Dir(/* what makes sense here? */),
		interp.Env(expand.ListEnviron(fmt.Sprintf("HOME=%s", homevar))), // TODO: configurable inclusion list.
		interp.ExecHandler(run.ExecHandler),
		interp.OpenHandler(run.OpenHandler),
		// interp.Params(), /* passed in by user */
		interp.StdIO(dn, dn, dn),
	)
	if err != nil {
		return err
	}

	err = int.Run(ctx, f)
	if err != nil {
		return err
	}

	fs := v.Manifest()
	for i := range fs {
		fs[i].Name = strings.ReplaceAll(fs[i].Name, homevar, "$HOME")
	}

	mfs := make([]*db.File, 0, len(fs))
	for _, f := range fs {
		mfs = append(mfs, &db.File{
			Name: f.Name,
			Dir:  f.Dir,
		})
	}

	m := &db.Manifest{
		URL:   url,
		Files: mfs,
	}

	d := db.DB{Root: filepath.Join(xdg.DataHome, "sumdog", "installed")}
	txn, err := d.Begin(m)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	if err := fileSync(fs); err != nil {
		return err
	}

	return txn.Commit()
}

type runner struct {
	builtin map[string]builtin.BuiltinFunc
	sham    map[string]builtin.BuiltinFunc
	vfs     *vfs.VFS

	permittedExec func(args []string) bool
	log           func(tag string, msg ...string)
}

func (r *runner) Log(tag string, msg ...string)    { r.log(tag, msg...) }
func (r *runner) ChDir(path string)                { r.vfs.ChDir(path) }
func (r *runner) File(path string) []byte          { return r.vfs.File(path) }
func (r *runner) Write(path string) io.WriteCloser { return r.vfs.Write(path) }
func (r *runner) MkDir(path string)                { r.vfs.MkDir(path) }
func (r *runner) Remove(path string)               { r.vfs.Remove(path) }
func (r *runner) Move(from, to string) error       { return r.vfs.Move(from, to) }

func (r *runner) ExecHandler(ctx context.Context, args []string) error {
	b, ok := r.builtin[args[0]]
	if !ok {
		shamcmd := strings.Join(args, " ")
		b, ok = r.sham[shamcmd]
		if !ok {
			if r.permittedExec(args) {
				return interp.DefaultExecHandler(2)(ctx, args)
			}

			return fmt.Errorf("unimplemented command: %s", args[0])
		}
	}

	hc := interp.HandlerCtx(ctx)
	return b(ctx, r, builtin.IOs{In: hc.Stdin, Out: hc.Stdout}, args[1:])
}

func (r *runner) OpenHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return devnull.New(), nil
	}

	return nil, fmt.Errorf("shell file opening not implemented")
}
