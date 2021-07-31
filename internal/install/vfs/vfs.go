package vfs

import (
	"errors"
	"io"
	"path/filepath"
)

type VFS struct {
	dir string // working directory

	files map[string]*File
}

type File struct {
	Dir bool
	b   []byte
}

// XXX: temporary
func (v *VFS) Manifest() []string {
	var out []string
	for f := range v.files {
		out = append(out, f)
	}

	return out
}

func (f *File) Write(p []byte) (int, error) {
	f.b = append(f.b, p...)
	return len(p), nil
}
func (f *File) Close() error { return nil }

func (v *VFS) ChDir(path string) {
	v.dir = path
}

func (v *VFS) MkDir(path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.dir, path)
	}

	if v.files == nil {
		v.files = make(map[string]*File)
	}

	v.files[path] = &File{Dir: true}
}

func (v *VFS) Remove(path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.dir, path)
	}

	if v.files == nil {
		return
	}

	// TODO: recurse
	delete(v.files, path)
}

// TODO: not a good enough return
func (v *VFS) File(path string) []byte {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.dir, path)
	}

	return v.files[path].b
}

func (v *VFS) Write(name string) io.WriteCloser {
	if !filepath.IsAbs(name) {
		name = filepath.Join(v.dir, name)
	}

	if v.files == nil {
		v.files = make(map[string]*File)
	}

	v.files[name] = &File{}

	return v.files[name]
}

func (v *VFS) Move(from, to string) error {
	// TODO: recurse

	if !filepath.IsAbs(from) {
		from = filepath.Join(v.dir, from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(v.dir, to)
	}

	f, ok := v.files[from]
	if !ok {
		return errors.New("file not found")
	}

	delete(v.files, from)

	t, ok := v.files[to]
	if !ok || !t.Dir {
		// TODO: check real FS.
		v.files[to] = f
	} else {
		v.files[filepath.Join(to, filepath.Base(from))] = f
	}

	return nil
}
