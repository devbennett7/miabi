// Copyright (c) 2026 Jonas Kaninda.
//
// This file is part of Miabi.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package stackcmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type quietUI struct{}

func (quietUI) Printf(string, ...any)  {}
func (quietUI) Info(string, ...any)    {}
func (quietUI) Success(string, ...any) {}
func (quietUI) Warn(string, ...any)    {}
func (quietUI) Confirm(string) bool    { return true }

func TestParseBackupName(t *testing.T) {
	for name, want := range map[string]DBBackup{
		"miabi-1.8.0-20261017T120000Z.dump":                  {Version: "1.8.0"},
		"miabi-1.8.0-rc.1-20261017T120000Z-pre-upgrade.dump": {Version: "1.8.0-rc.1", Label: "pre-upgrade"},
		"miabi-latest-20261017T120000Z-pre-restore.dump":     {Version: "latest", Label: "pre-restore"},
	} {
		got, ok := parseBackupName(name)
		if !ok {
			t.Errorf("%s: not parsed", name)
			continue
		}
		if got.Version != want.Version || got.Label != want.Label {
			t.Errorf("%s: got version=%q label=%q", name, got.Version, got.Label)
		}
		if !got.Taken.Equal(time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)) {
			t.Errorf("%s: taken = %v", name, got.Taken)
		}
	}
	for _, name := range []string{"miabi.dump", "miabi-1.8.0-20261017T120000Z.dump.partial", "notes.txt"} {
		if _, ok := parseBackupName(name); ok {
			t.Errorf("%s: parsed, want ignored", name)
		}
	}
}

func TestNewerThan(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.9.0", "1.8.5", true},
		{"1.10.0", "1.9.0", true},
		{"1.8.0", "1.8.0", false},
		{"1.8.0", "1.9.0", false},
		{"v1.9.0", "1.8.0", true},
		{"1.9.0-rc.1", "1.8.0", true},
		{"latest", "1.8.0", false},
		{"1.9.0", "digest", false},
	}
	for _, c := range cases {
		if got := newerThan(c.a, c.b); got != c.want {
			t.Errorf("newerThan(%q, %q) = %v", c.a, c.b, got)
		}
	}
}

func TestPruneAutoKeepsManualDumps(t *testing.T) {
	dir := t.TempDir()
	touch := func(name string) {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for d := 1; d <= autoKeep+2; d++ {
		touch(time.Date(2026, 10, d, 0, 0, 0, 0, time.UTC).Format("miabi-1.8.0-" + backupTimeLayout + "-pre-restore.dump"))
	}
	touch("miabi-1.7.0-20200101T000000Z.dump")
	touch("miabi-1.7.0-20200101T000000Z-before-upgrade.dump")

	pruneAuto(dir, LabelPreRestore, quietUI{})

	all, err := ListBackups(dir)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{}
	for _, b := range all {
		counts[b.Label]++
	}
	if counts[LabelPreRestore] != autoKeep || counts[""] != 1 || counts["before-upgrade"] != 1 {
		t.Fatalf("after prune: %v", counts)
	}
	if all[0].Taken.Day() != autoKeep+2 {
		t.Errorf("newest pre-restore dump was pruned: %v", all[0].Path)
	}
}

func TestResolveBackupByBareName(t *testing.T) {
	manifest := filepath.Join(t.TempDir(), "miabi.yaml")
	dir := BackupDir(manifest)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	name := "miabi-1.8.0-20261017T120000Z.dump"
	if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := resolveBackup(manifest, name)
	if err != nil || got != filepath.Join(dir, name) {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := resolveBackup(manifest, "missing.dump"); err == nil {
		t.Fatal("want an error for a missing dump")
	}
}
