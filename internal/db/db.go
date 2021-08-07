// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"encoding/base32"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ManifestState string

const (
	UnknownManifestState       ManifestState = "unknown"
	InstalledManifestState     ManifestState = "installed"
	PendingManifestState       ManifestState = "pending"
	BrokenInstallManifestState ManifestState = "broken_install"
	BrokenManifestState        ManifestState = "broken"
	PendingDeleteManifestState ManifestState = "pending_delete"
	BrokenDeleteManifestState  ManifestState = "broken_delete"
)

type Manifest struct {
	// Serialized in this order.

	Name        string        `yaml:"name,omitempty"`
	Description string        `yaml:"description,omitempty"`
	State       ManifestState `yaml:"state"`

	URL string `yaml:"url"` // TODO: checksum
	// TODO: secondary URLs + checksum

	// TODO: include values of env vars when installed.

	// TODO: include created symlinks

	Files []*File `yaml:"files"`
}

type File struct {
	Name string `yaml:"name"`
	Dir  bool   `yaml:"dir,omitempty"`

	// TODO: mode, checksum
}

type DB struct {
	Root string
}

func (db *DB) Begin(m *Manifest) (*Txn, error) {
	if err := os.MkdirAll(db.Root, 0700); err != nil {
		return nil, err
	}

	m.State = PendingManifestState
	t := &Txn{db: db, m: m, fail: BrokenInstallManifestState, commit: func(t *Txn) error {
		t.m.State = InstalledManifestState
		return t.writeManifest()
	}}

	if err := t.writeManifest(); err != nil {
		return nil, err
	}

	return t, nil
}

func (db *DB) Delete(m *Manifest) (*Txn, error) {
	m.State = PendingManifestState
	t := &Txn{db: db, m: m, fail: BrokenDeleteManifestState, commit: func(t *Txn) error {
		if err := os.Remove(fname(t.db.Root, t.m.URL)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		return nil
	}}

	if err := t.writeManifest(); err != nil {
		return nil, err
	}

	return t, nil
}

type Txn struct {
	db   *DB
	m    *Manifest
	done bool

	fail   ManifestState
	commit func(t *Txn) error
}

func (t *Txn) Rollback() {
	if t.done {
		return
	}

	t.m.State = t.fail
	// TODO: panic instead?
	_ = t.writeManifest()

}

func (t *Txn) Commit() error {
	err := t.commit(t)
	if err == nil {
		t.done = true
	}

	return err
}

func (t *Txn) writeManifest() error {
	d, err := yaml.Marshal(t.m)
	if err != nil {
		return err
	}

	// base 32 hex encoding preserves alphabetic ordering
	return ioutil.WriteFile(fname(t.db.Root, t.m.URL), d, 0600)
}

func (d *DB) Query(url string) (*Manifest, error) {
	return openManifest(fname(d.Root, url))
}

func (d *DB) List(fn func(m *Manifest) error) error {
	des, err := os.ReadDir(d.Root)
	if err != nil {
		return err
	}

	for _, de := range des {
		if de.IsDir() {
			continue
		}

		if filepath.Ext(de.Name()) != ".yaml" {
			continue
		}

		m, err := openManifest(filepath.Join(d.Root, de.Name()))
		if err != nil {
			return err
		}

		if err := fn(m); err != nil {
			return err
		}
	}
	return nil
}

func openManifest(path string) (*Manifest, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func fname(root, url string) string {
	return filepath.Join(root, base32.HexEncoding.EncodeToString([]byte(url))+".yaml")
}
