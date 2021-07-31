package install

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/jbowes/sumdog/internal/install/vfs"
)

func fileSync(fs []*vfs.ManifestEntry) error {
	// TODO: attempt to cleanup on error? or leave for broken.
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
		f, err := os.OpenFile(name, os.O_WRONLY, 0700)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, bytes.NewReader(e.Data))
		if err != nil {
			return err
		}

		// TODO: support windows
		bin := filepath.Join(xdg.Home, ".local", "bin")

		rel, err := filepath.Rel(bin, name)
		if err != nil {
			rel = name
		}

		if err := os.MkdirAll(bin, 0700); err != nil {
			return err
		}

		// TODO: if exec only
		return os.Symlink(rel, filepath.Join(bin, filepath.Base(name)))
	}

	return nil
}
