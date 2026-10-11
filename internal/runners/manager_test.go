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

package runners

import (
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

// fakeState records MarkConnected/MarkDisconnected calls so the manager's
// status wiring can be asserted without a database.
type fakeState struct {
	connected   map[uint]string // id -> last version reported
	disconnects map[uint]int
}

func newFakeState() *fakeState {
	return &fakeState{connected: map[uint]string{}, disconnects: map[uint]int{}}
}

func (f *fakeState) MarkConnected(id uint, facts models.RunnerFacts) { f.connected[id] = facts.Version }
func (f *fakeState) MarkDisconnected(id uint)                        { f.disconnects[id]++ }

// A freshly-built manager reports every runner offline and hands back no session.
func TestManagerOfflineByDefault(t *testing.T) {
	m := NewManager(newFakeState())
	if m.Connected(7) {
		t.Error("runner 7 should be offline with no tunnel")
	}
	if _, ok := m.Session(7); ok {
		t.Error("no session should exist for an unconnected runner")
	}
}

// Disconnecting an unknown runner is a safe no-op (it never panics on a missing
// session), so a delete/disable of an offline runner is always safe to call.
func TestDisconnectUnknownIsNoop(t *testing.T) {
	m := NewManager(newFakeState())
	m.Disconnect(999) // must not panic
}
