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

// Package webhooksig verifies inbound provider push webhooks.
package webhooksig

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// MaxBodyBytes caps an inbound webhook body. These endpoints are unauthenticated
// until the signature is checked, and the signature cannot be checked without
// first reading the body, so the read itself has to be bounded. GitHub caps its
// own payloads at 25 MB; a push event carries at most 20 commits and runs well
// under this.
const MaxBodyBytes int64 = 1 << 20

// SignatureHeader is GitHub's HMAC header; TokenHeader is GitLab's bare-secret one.
const (
	SignatureHeader = "X-Hub-Signature-256"
	TokenHeader     = "X-Gitlab-Token"
)

// Verify reports whether signature authenticates body under secret, accepting
// either GitHub's X-Hub-Signature-256 HMAC-SHA256 scheme or GitLab's
// X-Gitlab-Token bare secret.
//
// An empty secret never verifies. Without that guard an unconfigured webhook is
// open to anyone: the HMAC of any body under an empty key is computable, so the
// attacker signs their own payload and it checks out.
func Verify(secret, signature string, body []byte) bool {
	signature = strings.TrimSpace(signature)
	if secret == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if hmac.Equal([]byte(expected), []byte(signature)) {
		return true
	}
	// GitLab sends the secret itself rather than an HMAC, so this compares one
	// secret against another and must not short-circuit on the first byte.
	return subtle.ConstantTimeCompare([]byte(signature), []byte(secret)) == 1
}
