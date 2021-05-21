package store

import (
	"fmt"
	"io"
	"path/filepath"
)

type Store struct {
	dir string // working directory

	files map[string]*file
}

type file struct {
	dir bool
	b   []byte
}

// XXX: temporary
func (s *Store) Manifest() {
	for f, _ := range s.files {
		fmt.Println(f)
	}
}

func (f *file) Write(p []byte) (int, error) {
	f.b = append(f.b, p...)
	return len(p), nil
}
func (f *file) Close() error { return nil }

func (s *Store) ChDir(path string) {
	s.dir = path
}

func (s *Store) MkDir(path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.dir, path)
	}

	if s.files == nil {
		s.files = make(map[string]*file)
	}

	s.files[path] = &file{dir: true}
}

func (s *Store) Remove(path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.dir, path)
	}

	if s.files == nil {
		return
	}

	// TODO: recurse
	delete(s.files, path)
}

// TODO: not a good enough return
func (s *Store) File(path string) []byte {
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.dir, path)
	}

	return s.files[path].b
}

func (s *Store) Write(name string) io.WriteCloser {
	if !filepath.IsAbs(name) {
		name = filepath.Join(s.dir, name)
	}

	if s.files == nil {
		s.files = make(map[string]*file)
	}

	s.files[name] = &file{}

	return s.files[name]
}
