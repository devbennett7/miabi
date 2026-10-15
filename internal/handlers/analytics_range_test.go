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
	"time"
)

func TestResolveWindow(t *testing.T) {
	now := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	day := 24 * time.Hour
	max := 90 * day
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }

	cases := []struct {
		name        string
		raw         string
		since, till time.Time
		wantErr     bool
	}{
		{name: "default", raw: "", since: now.Add(-day), till: now},
		{name: "relative", raw: "60d", since: now.Add(-60 * day), till: now},
		{name: "relative capped", raw: "365d", since: now.Add(-max), till: now},
		{name: "dates", raw: "2026-09-01..2026-09-16", since: at("2026-09-01T00:00:00Z"), till: at("2026-09-16T00:00:00Z")},
		{name: "rfc3339 with offset", raw: "2026-09-01T00:00:00+02:00..2026-09-02T00:00:00+02:00", since: at("2026-08-31T22:00:00Z"), till: at("2026-09-01T22:00:00Z")},
		{name: "end clamped to now", raw: "2026-10-14..2026-10-20", since: at("2026-10-14T00:00:00Z"), till: now},
		{name: "start clamped to retention", raw: "2026-01-01..2026-10-01", since: now.Add(-max), till: at("2026-10-01T00:00:00Z")},
		{name: "reversed", raw: "2026-09-10..2026-09-01", wantErr: true},
		{name: "entirely past retention", raw: "2026-01-01..2026-02-01", wantErr: true},
		{name: "garbage bound", raw: "yesterday..2026-09-01", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			since, until, err := resolveWindow(tc.raw, now, max)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got [%s, %s)", since, until)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !since.Equal(tc.since) || !until.Equal(tc.till) {
				t.Fatalf("got [%s, %s), want [%s, %s)", since, until, tc.since, tc.till)
			}
		})
	}
}
