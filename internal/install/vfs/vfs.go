// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type VFS struct {
	Dir   func(context.Context) string // Get the working directory
	files map[string]*file
}

type file struct {
	dir bool // TODO: could track -p here too, inferring other dirs created, helping with cleanup
	b   []byte
}

// fs.Fileinfo interface.
func (f file) Name() string       { return "" } // TODO: implement me
func (f file) Size() int64        { return int64(len(f.b)) }
func (f file) Mode() fs.FileMode  { return 0 }          // TODO: implement me
func (f file) ModTime() time.Time { return time.Now() } // TODO: implement me
func (f file) IsDir() bool        { return f.dir }
func (f file) Sys() interface{}   { return nil }

func (f *file) Write(p []byte) (int, error) {
	f.b = append(f.b, p...)
	return len(p), nil
}
func (f *file) Close() error { return nil }

func (v *VFS) Stat(ctx context.Context, path string) (fs.FileInfo, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.Dir(ctx), path)
	} else {
		path = filepath.Clean(path)
	}

	if _, ok := v.files[path]; !ok {
		return nil, fs.ErrNotExist
	}

	return v.files[path], nil
}

func (v *VFS) MkDir(ctx context.Context, path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.Dir(ctx), path)
	} else {
		path = filepath.Clean(path)
	}

	if v.files == nil {
		v.files = make(map[string]*file)
	}

	v.files[path] = &file{dir: true}
}

func (v *VFS) Remove(ctx context.Context, path string) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.Dir(ctx), path)
	} else {
		path = filepath.Clean(path)
	}

	if v.files == nil {
		return
	}

	for k := range v.files {
		// TODO: Not the best prefix comparision, WRT case-insensitivity
		if k == path || strings.HasPrefix(k, path) && len(k) > len(path) && k[len(path)] == filepath.Separator {
			delete(v.files, k)
		}
	}
}

// TODO: not a good enough return
func (v *VFS) File(ctx context.Context, path string) []byte {
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.Dir(ctx), path)
	} else {
		path = filepath.Clean(path)
	}

	return v.files[path].b
}

func (v *VFS) Write(ctx context.Context, name string) io.WriteCloser {
	if !filepath.IsAbs(name) {
		name = filepath.Join(v.Dir(ctx), name)
	} else {
		name = filepath.Clean(name)
	}

	if v.files == nil {
		v.files = make(map[string]*file)
	}

	v.files[name] = &file{}

	return v.files[name]
}

func (v *VFS) Move(ctx context.Context, from, to string) error {
	if !filepath.IsAbs(from) {
		from = filepath.Join(v.Dir(ctx), from)
	} else {
		from = filepath.Clean(from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(v.Dir(ctx), to)
	} else {
		to = filepath.Clean(to)
	}

	f, ok := v.files[from]
	if !ok {
		return errors.New("file not found")
	}

	// Figure out if we're moving into a directory or not.
	// TODO: this logic could be in the move command
	t, ok := v.files[to]
	if ok && t.dir {
		// TODO: check real FS.
		to = filepath.Join(to, filepath.Base(from))
	}

	delete(v.files, from)
	v.files[to] = f

	for k, f := range v.files {
		// TODO: Not the best prefix comparision, WRT case-insensitivity
		if strings.HasPrefix(k, from) && len(k) > len(from) && k[len(from)] == filepath.Separator {
			newTo := strings.Replace(k, from, to, 1)
			delete(v.files, k)
			v.files[newTo] = f
		}
	}

	return nil
}

// TODO: use stat info here
type ManifestEntry struct {
	Name string
	Dir  bool

	Data []byte
}

// XXX: temporary
func (v *VFS) Manifest() []*ManifestEntry {

	var out []*ManifestEntry
	for n, f := range v.files {
		out = append(out, &ManifestEntry{Name: n, Dir: f.dir, Data: f.b})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.Compare(out[i].Name, out[j].Name) < 1
	})

	return out
}
