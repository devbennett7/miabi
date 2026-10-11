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

package wsbundle

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestEncodeWithoutPassphraseWritesPlainJSON(t *testing.T) {
	body, err := Encode(testState(), "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(body), "{") || !strings.Contains(string(body), "sk_live_supersecret") {
		t.Fatalf("an unencrypted state file should be readable JSON, got %.40q", body)
	}
	out, err := Decode(body, "")
	if err != nil {
		t.Fatal(err)
	}
	if out.Workspace.Name != "shop" || out.Secrets[0].Value != "sk_live_supersecret" {
		t.Fatalf("round trip lost state: %+v", out)
	}
}

// A workspace that sets a passphrase after making plain bundles must still restore them.
func TestDecodeIgnoresThePassphraseForAPlainFile(t *testing.T) {
	body, _ := Encode(testState(), "")
	if _, err := Decode(body, goodPass); err != nil {
		t.Fatalf("decode plain with a passphrase set: %v", err)
	}
}

func TestDecodeSealedFile(t *testing.T) {
	sealed, err := Encode(testState(), goodPass)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Decode(sealed, goodPass); err != nil {
		t.Fatalf("decode sealed: %v", err)
	}
	// Encryption removed after the bundle was made: say what is needed, not "corrupt file".
	if _, err := Decode(sealed, ""); !errors.Is(err, ErrPassphraseRequired) {
		t.Fatalf("err = %v, want ErrPassphraseRequired", err)
	}
}

func TestDecodeRejectsGarbage(t *testing.T) {
	if _, err := Decode([]byte("not a bundle"), ""); !errors.Is(err, ErrNotState) {
		t.Fatalf("err = %v, want ErrNotState", err)
	}
}

func TestInfoNoticeFollowsEncryption(t *testing.T) {
	for _, enc := range []bool{true, false} {
		i := &Info{Ref: NewRef("shop", time.Unix(1_760_000_000, 0).UTC()), Encrypted: enc}
		body, err := EncodeInfo(i)
		if err != nil {
			t.Fatal(err)
		}
		plain := strings.Contains(string(body), "NOT encrypted")
		if plain == enc {
			t.Errorf("encrypted=%v: notice = %q", enc, i.Notice)
		}
	}
}
