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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/miabi-io/miabi/pkg/stack/docker"
)

// DBBackupExt is the extension of a platform database dump: pg_dump's custom format, which is
// compressed and lets a restore load into a fresh database.
const DBBackupExt = ".dump"

const (
	dbName        = "miabi"
	dbUser        = "miabi"
	dbRestoreName = "miabi_restore"
	dumpMount     = "/backup"
)

// DumpDatabase writes a pg_dump of the platform database to dir/name. It runs the client from the
// manifest's Postgres image, so pg_dump is never older than the server it reads.
//
// The dump lands under a temporary name and is renamed only once pg_dump exits cleanly, so a killed
// or failed run never leaves a file that looks like a backup.
func (s *Service) DumpDatabase(ctx context.Context, m *Manifest, dir, name string) (string, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("create %s: %w", dir, err)
	}
	final := filepath.Join(dir, name)
	partial := final + ".partial"
	_ = os.Remove(partial)

	script := fmt.Sprintf("pg_dump -h %s -U %s -d %s -Fc -f %s",
		ContainerPostgres, dbUser, dbName, shQuote(dumpMount+"/"+filepath.Base(partial)))
	if err := s.runPGClient(ctx, m, "db-dump", dir, false, script); err != nil {
		_ = os.Remove(partial)
		return "", fmt.Errorf("pg_dump: %w", err)
	}
	if err := os.Chmod(partial, 0o600); err != nil {
		_ = os.Remove(partial)
		return "", err
	}
	if err := os.Rename(partial, final); err != nil {
		_ = os.Remove(partial)
		return "", err
	}
	return final, nil
}

// RestoreDatabase replaces the platform database with the dump at file. The control plane must
// already be stopped: it would otherwise reconnect between the drop and the rename.
//
// The dump is loaded into a scratch database first and swapped in only when pg_restore succeeds, so
// a corrupt or truncated file leaves the live database untouched. Restoring into a fresh database
// rather than with --clean also drops tables a newer migration added after the dump was taken.
func (s *Service) RestoreDatabase(ctx context.Context, m *Manifest, file string) error {
	if _, err := os.Stat(file); err != nil {
		return err
	}
	psql := fmt.Sprintf("psql -q -v ON_ERROR_STOP=1 -h %s -U %s -d postgres", ContainerPostgres, dbUser)
	script := strings.Join([]string{
		"set -e",
		fmt.Sprintf("%s -c 'DROP DATABASE IF EXISTS %s WITH (FORCE)'", psql, dbRestoreName),
		fmt.Sprintf("createdb -h %s -U %s -O %s %s", ContainerPostgres, dbUser, dbUser, dbRestoreName),
		fmt.Sprintf("pg_restore -h %s -U %s -d %s --no-owner --exit-on-error --single-transaction %s || { %s -c 'DROP DATABASE %s'; exit 1; }",
			ContainerPostgres, dbUser, dbRestoreName, shQuote(dumpMount+"/"+filepath.Base(file)), psql, dbRestoreName),
		fmt.Sprintf("%s -c 'DROP DATABASE %s WITH (FORCE)'", psql, dbName),
		fmt.Sprintf("%s -c 'ALTER DATABASE %s RENAME TO %s'", psql, dbRestoreName, dbName),
	}, "\n")
	if err := s.runPGClient(ctx, m, "db-restore", filepath.Dir(file), true, script); err != nil {
		return fmt.Errorf("restore: %w", err)
	}
	return nil
}

// StopControlPlane removes the control-plane container so nothing writes to the database while it is
// replaced. It is stateless — Converge recreates it from the manifest.
func (s *Service) StopControlPlane(ctx context.Context) error {
	if err := s.dc.RemoveContainer(ctx, ContainerControlPlane, true); err != nil && !isNotFound(err) {
		return fmt.Errorf("stop %s: %w", ContainerControlPlane, err)
	}
	return nil
}

// runPGClient runs script in a throwaway container from the Postgres image, on the database's own
// network and with dir bound at /backup. It runs as the calling user so the dump is owned by
// whoever asked for it, not by root.
func (s *Service) runPGClient(ctx context.Context, m *Manifest, role, dir string, readOnly bool, script string) error {
	pg, err := s.dc.InspectContainer(ctx, ContainerPostgres)
	if err != nil {
		if isNotFound(err) {
			return fmt.Errorf("%s is not running — start the stack first (miabi setup)", ContainerPostgres)
		}
		return err
	}
	if !strings.EqualFold(pg.State, "running") && pg.State != "" {
		return fmt.Errorf("%s is %s — start the stack first (miabi setup)", ContainerPostgres, pg.State)
	}
	// Join the network Postgres is actually on: an install predating the internal network still has
	// it on the proxy network until its next converge.
	network := m.InternalNetwork.Name
	if len(pg.Networks) > 0 && !onNetwork(pg, network) {
		network = pg.Networks[0].Name
	}

	host, err := s.requireHostPath(ctx, dir)
	if err != nil {
		return err
	}
	if err := ensureImage(ctx, s.dc, m.Images.Postgres, s.log); err != nil {
		return err
	}
	spec := docker.RunSpec{
		Name:       "mb-" + role,
		Image:      m.Images.Postgres,
		Entrypoint: []string{"sh", "-c"},
		Cmd:        []string{script},
		// HOME: the calling uid has no passwd entry in the image, and psql looks for ~/.psqlrc.
		Env:      []string{"PGPASSWORD=" + m.Secrets.DBPassword, "HOME=/tmp", "PGOPTIONS=-c client_min_messages=warning"},
		User:     fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		Networks: []string{network},
		Binds:    []docker.BindMount{{Source: host, Target: dumpMount, ReadOnly: readOnly}},
		Labels:   map[string]string{docker.LabelPartOf: docker.PartOfMiabi, docker.LabelRole: role},
	}
	_ = s.dc.RemoveContainer(ctx, spec.Name, true) // a leftover from a killed run

	code, out, err := s.dc.RunOneShot(ctx, spec)
	if err != nil {
		return err
	}
	if code != 0 {
		msg := strings.TrimSpace(out)
		if msg == "" {
			msg = fmt.Sprintf("exit code %d", code)
		}
		return errors.New(msg)
	}
	return nil
}

func onNetwork(c docker.Container, name string) bool {
	for _, n := range c.Networks {
		if n.Name == name {
			return true
		}
	}
	return false
}

func shQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }
