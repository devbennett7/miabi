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
	"strings"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/crypto"
)

func TestCredentialURL(t *testing.T) {
	crypto.Init("test-master-key-for-gitrepo")
	tok, err := crypto.Encrypt("ghp_secrettoken")
	if err != nil {
		t.Fatal(err)
	}

	// Public / nil credential: URL unchanged (normalized to .git).
	if got, err := CredentialURL("https://github.com/acme/web", nil, nil); err != nil || got != "https://github.com/acme/web.git" {
		t.Errorf("public: got %q err %v", got, err)
	}

	// Token: credential embedded as https://user:token@host/….
	g := &models.GitRepository{AuthType: models.GitAuthToken, Username: "x-access-token", Secret: tok}
	got, err := CredentialURL("https://github.com/acme/web.git", g, nil)
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if got != "https://x-access-token:ghp_secrettoken@github.com/acme/web.git" {
		t.Errorf("token URL = %q", got)
	}

	// Empty username defaults to the provider-agnostic x-access-token.
	g2 := &models.GitRepository{AuthType: models.GitAuthToken, Secret: tok}
	if got2, _ := CredentialURL("https://github.com/acme/web.git", g2, nil); !strings.Contains(got2, "x-access-token:ghp_secrettoken@") {
		t.Errorf("default username not applied: %q", got2)
	}

	// SSH key: can't be embedded in a URL — clear error.
	gs := &models.GitRepository{AuthType: models.GitAuthSSH, Secret: tok}
	if _, err := CredentialURL("git@github.com:acme/web.git", gs, nil); err != ErrSSHUnsupportedOnRunner {
		t.Errorf("ssh: want ErrSSHUnsupportedOnRunner, got %v", err)
	}
}
