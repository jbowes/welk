package install

import (
	"bytes"
	"io"
	"os"

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
		f, err := os.OpenFile(name, os.O_WRONLY, 0700)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, bytes.NewReader(e.Data))
		if err != nil {
			return err
		}
	}

	return nil
}
