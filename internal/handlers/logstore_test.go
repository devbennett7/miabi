// Copyright (c) 2026 Jonas Kaninda.
//
// This file is part of Miabi.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"testing"

	"github.com/miabi-io/miabi/internal/logstore"
)

func TestReplayLogHistory(t *testing.T) {
	be, err := logstore.NewFSBackend(t.TempDir())
	if err != nil {
		t.Fatalf("backend: %v", err)
	}
	store := logstore.New(be, logstore.Config{Compress: true})
	ref := logstore.DeploymentRef(1, 2, 3)
	if _, err := store.Externalize(ref, "alpha\nbravo\ncharlie\n"); err != nil {
		t.Fatalf("externalize: %v", err)
	}

	// With a valid ref the full history comes from the store, not the tail.
	got := replayLogHistory(store, ref, "only-tail\n")
	if len(got) != 3 || got[0] != "alpha" || got[2] != "charlie" {
		t.Errorf("store history = %v, want [alpha bravo charlie]", got)
	}

	// An empty ref (in-progress / pre-migration) replays the DB tail.
	got = replayLogHistory(store, "", "tail-one\ntail-two\n")
	if len(got) != 2 || got[1] != "tail-two" {
		t.Errorf("tail replay = %v, want [tail-one tail-two]", got)
	}

	// A dangling ref (object swept) falls back to the tail rather than erroring.
	got = replayLogHistory(store, logstore.DeploymentRef(9, 9, 9), "fallback\n")
	if len(got) != 1 || got[0] != "fallback" {
		t.Errorf("missing-object fallback = %v, want [fallback]", got)
	}

	// A disabled (nil) store always uses the tail.
	got = replayLogHistory(nil, ref, "disabled-tail\n")
	if len(got) != 1 || got[0] != "disabled-tail" {
		t.Errorf("disabled store = %v, want [disabled-tail]", got)
	}
}
