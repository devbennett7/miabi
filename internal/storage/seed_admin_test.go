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

package storage

import (
	"errors"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func adminDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatal(err)
	}
	return db
}

var errBadCredentials = errors.New("MIABI_ADMIN_PASSWORD must be set")

func TestSeedAdminRefusesInvalidCredentialsOnFirstBoot(t *testing.T) {
	db := adminDB(t)
	_, err := SeedAdmin(db, "admin@acme.test", "admin@1234", func() error { return errBadCredentials })
	if !errors.Is(err, ErrAdminSeedInvalid) {
		t.Fatalf("err = %v, want ErrAdminSeedInvalid", err)
	}
	var n int64
	db.Model(&models.User{}).Count(&n)
	if n != 0 {
		t.Fatalf("%d users created from refused credentials", n)
	}
}

// Once an admin exists the credentials are never read, so they can be left unset.
func TestSeedAdminSkipsValidationWhenAnAdminExists(t *testing.T) {
	db := adminDB(t)
	first, err := SeedAdmin(db, "admin@acme.test", "a-real-password", nil)
	if err != nil {
		t.Fatal(err)
	}
	called := false
	again, err := SeedAdmin(db, "", "", func() error { called = true; return errBadCredentials })
	if err != nil || called || again.ID != first.ID {
		t.Fatalf("admin %v, err %v, validate called %v; want the existing admin without validating", again, err, called)
	}
}
