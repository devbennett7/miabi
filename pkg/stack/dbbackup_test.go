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

package stack

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/miabi-io/miabi/pkg/stack/docker"
)

// dumpEngine runs one-shots by calling run, so a test decides what pg_dump "did".
type dumpEngine struct {
	*fakeEngine
	specs []docker.RunSpec
	run   func(docker.RunSpec) (int, string)
}

func (d *dumpEngine) RunOneShot(_ context.Context, spec docker.RunSpec) (int, string, error) {
	d.specs = append(d.specs, spec)
	code, out := d.run(spec)
	return code, out, nil
}

func newDumpEngine(run func(docker.RunSpec) (int, string)) *dumpEngine {
	return &dumpEngine{
		fakeEngine: &fakeEngine{containers: map[string][]string{ContainerPostgres: {DefaultInternalNetwork}}},
		run:        run,
	}
}

func dumpManifest() *Manifest {
	m := Defaults("miabi/miabi:1.8.0")
	m.Domain = "miabi.example.com"
	m.Secrets.DBPassword = "s3cret"
	return m
}

func TestDumpDatabaseWritesOnlyACompleteFile(t *testing.T) {
	dir := t.TempDir()
	dc := newDumpEngine(func(spec docker.RunSpec) (int, string) {
		host := spec.Binds[0].Source
		_ = os.WriteFile(filepath.Join(host, "x.dump.partial"), []byte("PGDMP"), 0o644)
		return 0, ""
	})
	svc := New(dc, nil, filepath.Join(dir, "miabi.yaml"))

	got, err := svc.DumpDatabase(context.Background(), dumpManifest(), dir, "x.dump")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(dir, "x.dump") {
		t.Fatalf("path = %s", got)
	}
	st, err := os.Stat(got)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o, want 600: the dump holds every secret in the database", st.Mode().Perm())
	}
	spec := dc.specs[0]
	if spec.Image != DefaultPostgresImage {
		t.Errorf("image = %s, want the server's own so pg_dump is never older than it", spec.Image)
	}
	if spec.Networks[0] != DefaultInternalNetwork {
		t.Errorf("network = %v", spec.Networks)
	}
	if !strings.Contains(spec.Cmd[0], "-h "+ContainerPostgres) {
		t.Errorf("script = %q", spec.Cmd[0])
	}
}

func TestDumpDatabaseFailureLeavesNoFile(t *testing.T) {
	dir := t.TempDir()
	dc := newDumpEngine(func(spec docker.RunSpec) (int, string) {
		_ = os.WriteFile(filepath.Join(spec.Binds[0].Source, "x.dump.partial"), []byte("PG"), 0o644)
		return 1, "pg_dump: error: connection refused"
	})
	svc := New(dc, nil, filepath.Join(dir, "miabi.yaml"))

	_, err := svc.DumpDatabase(context.Background(), dumpManifest(), dir, "x.dump")
	if err == nil || !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("err = %v", err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("left %v behind; a failed dump must not look like a backup", entries)
	}
}

func TestDumpDatabaseNeedsPostgres(t *testing.T) {
	dc := newDumpEngine(nil)
	dc.containers = map[string][]string{}
	svc := New(dc, nil, "miabi.yaml")
	_, err := svc.DumpDatabase(context.Background(), dumpManifest(), t.TempDir(), "x.dump")
	if err == nil || !strings.Contains(err.Error(), "not running") {
		t.Fatalf("err = %v", err)
	}
}

func TestRestoreDatabaseSwapsOnlyAfterLoading(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "it's.dump")
	if err := os.WriteFile(file, []byte("PGDMP"), 0o600); err != nil {
		t.Fatal(err)
	}
	dc := newDumpEngine(func(docker.RunSpec) (int, string) { return 0, "" })
	svc := New(dc, nil, filepath.Join(dir, "miabi.yaml"))

	if err := svc.RestoreDatabase(context.Background(), dumpManifest(), file); err != nil {
		t.Fatal(err)
	}
	spec := dc.specs[0]
	if !spec.Binds[0].ReadOnly {
		t.Error("the dump directory should be bound read-only")
	}
	script := spec.Cmd[0]
	load := strings.Index(script, "pg_restore")
	drop := strings.Index(script, "DROP DATABASE "+dbName+" ")
	if load < 0 || drop < 0 || drop < load {
		t.Errorf("the live database must be dropped only after pg_restore succeeds:\n%s", script)
	}
	if !strings.Contains(script, `'/backup/it'\''s.dump'`) {
		t.Errorf("file name not shell-quoted:\n%s", script)
	}
}
