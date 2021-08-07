// Package filesync provides utilities to synchronize manifests with the file system.
package filesync

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/jbowes/welk/internal/db"
	"github.com/jbowes/welk/internal/install/vfs"
)

func Sync(fs []*vfs.ManifestEntry) error {
	// TODO: attempt to cleanup on error? or leave for broken.

	os.Setenv("XDG_DATA_HOME", xdg.DataHome)

	for _, e := range fs {

		// TODO: windows support needed.
		name := os.ExpandEnv(e.Name)

		if e.Dir {
			if err := os.MkdirAll(name, 0700); err != nil {
				return err
			}

			continue
		}

		// TODO: proper mode
		f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}

		_, err = io.Copy(f, bytes.NewReader(e.Data))
		if err != nil {
			return err
		}

		// TODO: support windows
		bin := filepath.Join(xdg.Home, ".local", "bin")
		if err := os.MkdirAll(bin, 0700); err != nil {
			return err
		}

		rel, err := filepath.Rel(bin, name)
		if err != nil {
			rel = name
		}

		// TODO: if exec only
		sym := filepath.Join(bin, filepath.Base(name))
		err = os.Symlink(rel, sym)
		switch {
		case err == nil:
		case errors.Is(err, os.ErrExist):
			s, sErr := filepath.EvalSymlinks(sym)
			if sErr != nil {
				return err
			}

			if s != name {
				// TODO: offer to replace symlink?
				return err
			}
		default:
			return err
		}

	}

	return nil
}

func Remove(fs []*db.File) error {
	// TODO: I don't like that Sync uses vfs and remove uses db.

	os.Setenv("XDG_DATA_HOME", xdg.DataHome)

	// TODO: report multiple errors.
	var seenErr error

	// go backwards, just so we clean up files in case of non-empty dirs.
	for i := len(fs) - 1; i >= 0; i-- {
		f := fs[i]

		// TODO: windows support needed.
		name := os.ExpandEnv(f.Name)

		if err := os.Remove(name); err != nil && !errors.Is(err, os.ErrNotExist) {
			seenErr = err
		}

		// TODO: support windows
		bin := filepath.Join(xdg.Home, ".local", "bin")

		// TODO: if exec only
		sym := filepath.Join(bin, filepath.Base(name))
		if err := os.Remove(sym); err != nil && !errors.Is(err, os.ErrNotExist) {
			seenErr = err
		}
	}

	return seenErr
}
