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

package cron

import "testing"

type fakeLeader struct{ leading bool }

func (l *fakeLeader) Leading() bool { return l.leading }

func TestTasksRunOnlyWhileLeadingUnlessLocal(t *testing.T) {
	m := NewManager(nil, nil, nil, nil, nil)
	l := &fakeLeader{}
	m.SetLeader(l)

	var shared, local int
	if err := m.RegisterTask("shared", 0, "Shared", "@every 1m", func() error { shared++; return nil }); err != nil {
		t.Fatal(err)
	}
	if err := m.RegisterLocalTask("local", 0, "Local", "@every 1m", func() error { local++; return nil }); err != nil {
		t.Fatal(err)
	}
	tick := func() {
		for _, key := range []string{taskKey("shared", 0), taskKey("local", 0)} {
			m.c.Entry(m.entries[key]).WrappedJob.Run()
		}
	}

	tick()
	if shared != 0 || local != 1 {
		t.Fatalf("standby tick ran shared=%d local=%d; want 0 and 1", shared, local)
	}
	for _, j := range m.Snapshot() {
		if j.Kind == "shared" && j.LastRunAt != nil {
			t.Fatal("a skipped task was recorded as run")
		}
	}

	l.leading = true
	tick()
	if shared != 1 || local != 2 {
		t.Fatalf("leader tick ran shared=%d local=%d; want 1 and 2", shared, local)
	}
}
