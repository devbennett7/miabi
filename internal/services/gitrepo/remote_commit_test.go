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

package gitrepo

import (
	"context"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

func TestMatchRemoteRef(t *testing.T) {
	const (
		mainSHA    = "1111111111111111111111111111111111111111"
		devSHA     = "2222222222222222222222222222222222222222"
		tagObject  = "3333333333333333333333333333333333333333"
		tagCommit  = "4444444444444444444444444444444444444444"
		lightTagSH = "5555555555555555555555555555555555555555"
	)
	refs := []*plumbing.Reference{
		plumbing.NewSymbolicReference(plumbing.HEAD, "refs/heads/main"),
		plumbing.NewHashReference("refs/heads/main", plumbing.NewHash(mainSHA)),
		plumbing.NewHashReference("refs/heads/dev", plumbing.NewHash(devSHA)),
		plumbing.NewHashReference("refs/tags/v1.0.0", plumbing.NewHash(tagObject)),
		plumbing.NewHashReference("refs/tags/v1.0.0^{}", plumbing.NewHash(tagCommit)),
		plumbing.NewHashReference("refs/tags/light", plumbing.NewHash(lightTagSH)),
	}
	cases := []struct {
		ref  string
		want string
		ok   bool
	}{
		{"", mainSHA, true},
		{"main", mainSHA, true},
		{"dev", devSHA, true},
		{"refs/heads/dev", devSHA, true},
		{"v1.0.0", tagCommit, true},
		{"refs/tags/v1.0.0", tagCommit, true},
		{"light", lightTagSH, true},
		{"missing", "", false},
	}
	for _, c := range cases {
		t.Run(c.ref, func(t *testing.T) {
			got, ok := matchRemoteRef(refs, c.ref)
			if got != c.want || ok != c.ok {
				t.Errorf("matchRemoteRef(%q) = (%q, %v), want (%q, %v)", c.ref, got, ok, c.want, c.ok)
			}
		})
	}
}

func TestRemoteCommitReturnsAFullHashWithoutListing(t *testing.T) {
	const sha = "1111111111111111111111111111111111111111"
	got, err := RemoteCommit(context.Background(), "https://invalid.example/repo.git", sha, nil)
	if err != nil || got != sha {
		t.Fatalf("RemoteCommit = (%q, %v), want (%q, nil)", got, err, sha)
	}
}
