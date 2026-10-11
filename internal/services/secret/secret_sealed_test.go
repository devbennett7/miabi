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

package secret

import (
	"errors"
	"testing"
)

func TestApplySealedCreatesThenOnlyBumpsOnARealChange(t *testing.T) {
	svc := multilineTestService(t)
	sec, err := svc.ApplySealed(1, "db", "pw-1", "fp-a", 1)
	if err != nil || sec.Version != 1 || sec.SealedFP != "fp-a" || sec.SealedKeyVersion != 1 {
		t.Fatalf("create = %+v, %v", sec, err)
	}

	sec, err = svc.ApplySealed(1, "db", "pw-1", "fp-b", 2)
	if err != nil || sec.Version != 1 || sec.SealedFP != "fp-b" || sec.SealedKeyVersion != 2 {
		t.Fatalf("re-seal of the same value must not bump the version: %+v, %v", sec, err)
	}

	sec, err = svc.ApplySealed(1, "db", "pw-2", "fp-c", 2)
	if err != nil || sec.Version != 2 {
		t.Fatalf("a new value must bump the version: %+v, %v", sec, err)
	}
	if got, _ := svc.Reveal(1, sec.ID); got != "pw-2" {
		t.Fatalf("stored value = %q", got)
	}
}

func TestHandSetValueClearsTheSealedMarker(t *testing.T) {
	svc := multilineTestService(t)
	sec, _ := svc.ApplySealed(1, "db", "pw-1", "fp-a", 1)
	sec, err := svc.Update(1, sec.ID, "typed-in-the-ui", "", nil)
	if err != nil || sec.SealedFP != "" || sec.SealedKeyVersion != 0 {
		t.Fatalf("Update = %+v, %v", sec, err)
	}
}

func TestApplySealedRefusesAManagedSecret(t *testing.T) {
	svc := multilineTestService(t)
	if _, err := svc.UpsertOwned(1, "database", 9, "db_url", "postgres://", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ApplySealed(1, "db_url", "x", "fp", 1); !errors.Is(err, ErrManaged) {
		t.Fatalf("want ErrManaged, got %v", err)
	}
}
