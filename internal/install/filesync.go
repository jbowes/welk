package install

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/jbowes/sumdog/internal/install/vfs"
)

func fileSync(fs []*vfs.ManifestEntry) error {
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
