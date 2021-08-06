package db

import (
	"encoding/base32"
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
	t := &Txn{db: db, m: m}

	if err := t.writeManifest(); err != nil {
		return nil, err
	}

	return t, nil
}

type Txn struct {
	db   *DB
	m    *Manifest
	done bool
}

func (t *Txn) Rollback() error {
	if t.done {
		return nil
	}

	t.m.State = BrokenInstallManifestState
	return t.writeManifest()

}

func (t *Txn) Commit() error {
	t.m.State = InstalledManifestState
	err := t.writeManifest()
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
