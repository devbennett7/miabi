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

package upgrade

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type rtAppFixture struct {
	ID          uint `gorm:"primaryKey"`
	RuntimeKind string
}

func (rtAppFixture) TableName() string { return "applications" }

type rtReleaseFixture struct {
	ID            uint `gorm:"primaryKey"`
	ApplicationID uint
	RuntimeKind   string
}

func (rtReleaseFixture) TableName() string { return "releases" }

// Existing releases take their app's runtime; one already recorded is left alone.
func TestReleaseRuntimeStep(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&rtAppFixture{}, &rtReleaseFixture{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&[]rtAppFixture{{ID: 1, RuntimeKind: "service"}, {ID: 2, RuntimeKind: "container"}, {ID: 3, RuntimeKind: ""}})
	db.Create(&[]rtReleaseFixture{
		{ID: 1, ApplicationID: 1}, {ID: 2, ApplicationID: 2}, {ID: 3, ApplicationID: 3},
		{ID: 4, ApplicationID: 1, RuntimeKind: "container"},
	})
	for range 2 {
		if err := releaseRuntimeStep(context.Background(), db); err != nil {
			t.Fatal(err)
		}
	}
	want := map[uint]string{1: "service", 2: "container", 3: "container", 4: "container"}
	var got []rtReleaseFixture
	db.Order("id").Find(&got)
	for _, r := range got {
		if r.RuntimeKind != want[r.ID] {
			t.Errorf("release %d runtime = %q, want %q", r.ID, r.RuntimeKind, want[r.ID])
		}
	}
}
