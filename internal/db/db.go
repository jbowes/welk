package db

import (
	"encoding/base64"
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
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	URL string `yaml:"url"` // TODO: checksum
	// TODO: secondary URLs + checksum

	Files []string `yaml:"files"` // TODO: dir or not, mode, checksum

	State ManifestState `yaml:"state"`
}

type DB struct {
	Root string
}

func (db *DB) Begin(m *Manifest) (*Txn, error) {
	if err := os.MkdirAll(db.Root, 0700); err != nil {
		return nil, err
	}

	m.State = PendingManifestState
	t := &Txn{db, m}

	if err := t.writeManifest(); err != nil {
		return nil, err
	}

	return t, nil
}

type Txn struct {
	db *DB
	m  *Manifest
}

func (t *Txn) Rollback() error {
	t.m.State = BrokenInstallManifestState
	return t.writeManifest()

}

func (t *Txn) Commit() error {
	t.m.State = InstalledManifestState
	return t.writeManifest()
}

func (t *Txn) writeManifest() error {
	d, err := yaml.Marshal(t.m)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(t.db.Root, base64.RawURLEncoding.EncodeToString([]byte(t.m.URL))+".yaml"), d, 0600)
}
