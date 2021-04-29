package store

import (
	"bytes"
	"io"
	"path/filepath"
)

type Store struct {
	dir string // working directory

	files map[string]file
}

type file struct {
	b bytes.Buffer
}

func (f *file) Write(p []byte) (int, error) { return f.b.Write(p) }
func (f *file) Close() error                { return nil }

func (s *Store) ChDir(path string) {
	s.dir = path
}

func (s *Store) Write(name string) io.WriteCloser {
	if !filepath.IsAbs(name) {
		name = filepath.Join(s.dir, name)
	}

	if s.files == nil {
		s.files = make(map[string]file)
	}

	s.files[name] = file{}

	f := s.files[name]
	return &f
}
