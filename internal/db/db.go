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

	Files []string `yaml:"files"` // TODO: dir or not, mode, checksum
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
	return ioutil.WriteFile(filepath.Join(t.db.Root, base32.HexEncoding.EncodeToString([]byte(t.m.URL))+".yaml"), d, 0600)
}
