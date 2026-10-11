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

package logintoken

import (
	"errors"
	"testing"
	"time"

	"github.com/miabi-io/miabi/internal/models"
)

func TestResolveScopes(t *testing.T) {
	cases := []struct {
		name    string
		in      []string
		want    []string
		wantErr error
	}{
		{"default is console-equivalent", nil, []string{models.ScopeRead, models.ScopeWrite, models.ScopeDeploy}, nil},
		{"empty defaults too", []string{}, []string{models.ScopeRead, models.ScopeWrite, models.ScopeDeploy}, nil},
		{"explicit read-only", []string{"read"}, []string{"read"}, nil},
		{"registry scopes allowed", []string{"registry_read"}, []string{"registry_read"}, nil},
		{"admin rejected", []string{"read", "admin"}, nil, ErrAdminScope},
		{"wildcard rejected", []string{"*"}, nil, ErrAdminScope},
		{"unknown rejected", []string{"nope"}, nil, ErrInvalidScope},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveScopes(tc.in)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("scopes = %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("scopes = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestResolveExpiry(t *testing.T) {
	s := &Service{ttl: 24 * time.Hour, maxTTL: 168 * time.Hour}
	hours := func(h int) *int { return &h }

	within := func(got, want time.Time) bool {
		d := got.Sub(want)
		return d > -time.Minute && d < time.Minute
	}

	if got := s.resolveExpiry(nil); !within(got, time.Now().Add(24*time.Hour)) {
		t.Errorf("default expiry = %v, want ~24h", time.Until(got))
	}
	if got := s.resolveExpiry(hours(1)); !within(got, time.Now().Add(1*time.Hour)) {
		t.Errorf("requested 1h expiry = %v", time.Until(got))
	}
	// Over the cap is clamped to maxTTL.
	if got := s.resolveExpiry(hours(1000)); !within(got, time.Now().Add(168*time.Hour)) {
		t.Errorf("over-cap expiry = %v, want ~168h", time.Until(got))
	}
	// Zero/negative falls back to the default.
	if got := s.resolveExpiry(hours(0)); !within(got, time.Now().Add(24*time.Hour)) {
		t.Errorf("zero-hours expiry = %v, want ~24h", time.Until(got))
	}
}
