package install

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/jbowes/sumdog/internal/db"
	"github.com/jbowes/sumdog/internal/forked/interp" // originally mvdan.cc/sh/v3/interp
	"github.com/jbowes/sumdog/internal/install/builtin"
	"github.com/jbowes/sumdog/internal/install/devnull"
	"github.com/jbowes/sumdog/internal/install/sham"
	"github.com/jbowes/sumdog/internal/install/vfs"
	"mvdan.cc/sh/v3/expand"
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

			if len(x.Args) > 0 && x.Args[0].Lit() == "command" {
				x.Args[0].Parts[0].(*syntax.Lit).Value = "sumdog-command"
			}

			// TODO: log cd, pushd, popd
		}
		return true
	})

	v := &vfs.VFS{}
	run := &runner{
		VFS: v,

		builtin:       builtin.Builtin,
		sham:          sham.Sham,
		permittedExec: permittedExec,
		log:           log,
	}

	// TODO: windows
	homevarU, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	homevar := "/" + homevarU.String()

	dn := devnull.New()

	int, err := interp.New(
		// interp.Dir(/* what makes sense here? */),
		interp.Env(expand.ListEnviron(fmt.Sprintf("HOME=%s", homevar))), // TODO: configurable inclusion list.
		interp.ExecHandler(run.ExecHandler),
		interp.OpenHandler(run.OpenHandler),
		interp.StatHandler(run.StatHandler),
		// interp.Params(), /* passed in by user */
		interp.StdIO(dn, dn, dn),
	)
	if err != nil {
		return err
	}

	pdU, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	pkgDir := "/" + pdU.String()

	int.Dir = pkgDir
	v.Dir = func(ctx context.Context) string { return interp.HandlerCtx(ctx).Dir }

	err = int.Run(ctx, f)
	if err != nil {
		return err
	}

	fmt.Println("Preparing to install")

	fs := v.Manifest()
	for i := range fs {
		fs[i].Name = strings.ReplaceAll(fs[i].Name, homevar, "$HOME")
		fmt.Println(fs[i].Name)
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
	*vfs.VFS

	builtin map[string]builtin.BuiltinFunc
	sham    map[string]builtin.BuiltinFunc

	permittedExec func(args []string) bool
	log           func(tag string, msg ...string)
}

func (r *runner) Log(tag string, msg ...string) { r.log(tag, msg...) }

var paths = []string{
	"/usr/bin",
	"/usr/local/bin",
}

func (r *runner) ExecHandler(ctx context.Context, args []string) error {
	cmd := args[0]
	if filepath.IsAbs(cmd) {
		found := false
		prefix := filepath.Dir(cmd)
		for _, p := range paths {
			if prefix == p {
				found = true
				break
			}
		}

		if found {
			cmd = filepath.Base(cmd)
		}
	}

	b, ok := r.builtin[cmd]
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
	err := b(ctx, r, builtin.IOs{In: hc.Stdin, Out: hc.Stdout}, args[1:])
	var e *builtin.ExitStatusError
	switch {
	case err == nil:
	case errors.As(err, &e):
		return interp.NewExitStatus(e.Status())
	default:
		return err
	}
	if err != nil {
		fmt.Printf("err in %s: %s\n", cmd, err)
	}

	return err
}

func (r *runner) OpenHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return devnull.New(), nil
	}

	// TODO: connect this to the VFS
	return nil, fmt.Errorf("shell file opening not implemented")
}

func (r *runner) StatHandler(ctx context.Context, path string) (os.FileInfo, error) {
	return r.Stat(ctx, path)
}
