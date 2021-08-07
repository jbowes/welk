// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sham

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jbowes/welk/internal/install/builtin"
	"golang.org/x/mod/semver"
)

const soloVersionParse = `python -c import sys; from distutils.version import LooseVersion; from json import loads as l; releases = l(sys.stdin.read()); releases = [release['tag_name'] for release in releases if not release['prerelease'] ];  releases.sort(key=LooseVersion, reverse=True); print('\n'.join(releases))`

// SoloVersionParse is used (at least) by Solo's wasme install script (https://run.solo.io/wasme/instal)
// to find the latest released version of the command.
// Versions are suppled to stdin as a json array.
// The output should be a newline delimited list of sorted release tags, excluding prereleases
func SoloVersionParse(ctx context.Context, host builtin.Host, ios builtin.IOs, args []string) error {
	host.Log("python") // TODO: need a better logging mechanism here

	var rels []struct {
		TagName    string `json:"tag_name"`
		Prerelease bool
	}

	dec := json.NewDecoder(ios.In)
	err := dec.Decode(&rels)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sort.Slice(rels, func(i, j int) bool { return semver.Compare(rels[i].TagName, rels[j].TagName) > 1 })

	for i, r := range rels {
		if !r.Prerelease {
			o := r.TagName
			if i != len(rels)-1 {
				o += "\n"
			}
			_, err := ios.Out.Write([]byte(o))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func init() { Sham[soloVersionParse] = SoloVersionParse }
