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

package repositories

import (
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type appRefRow struct {
	ID          uint `gorm:"primaryKey"`
	UID         string
	WorkspaceID uint
	Name        string
	DeletedAt   gorm.DeletedAt
}

func (appRefRow) TableName() string { return "applications" }

func TestIDByRef(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&appRefRow{}); err != nil {
		t.Fatal(err)
	}
	rows := []appRefRow{
		{ID: 1, UID: "6f1c2a3e-1111-4c4c-9a9a-000000000001", WorkspaceID: 7, Name: "api"},
		{ID: 2, UID: "6f1c2a3e-1111-4c4c-9a9a-000000000002", WorkspaceID: 7, Name: "2024"},
		{ID: 11, UID: "6f1c2a3e-1111-4c4c-9a9a-000000000011", WorkspaceID: 8, Name: "api"},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewApplicationRepository(db)
	cases := []struct {
		ws   uint
		ref  string
		want uint
	}{
		{7, "api", 1},
		{7, "API", 1},
		{8, "api", 11}, // same name, other workspace
		{7, "1", 1},    // id
		{7, "2024", 2}, // digits: no app with id 2024 here, so the name
		{7, "6f1c2a3e-1111-4c4c-9a9a-000000000001", 1}, // uid
		{7, "11", 0}, // an id from another workspace
		{7, "6f1c2a3e-1111-4c4c-9a9a-000000000011", 0}, // a uid from another workspace
		{7, "nope", 0},
		{7, "", 0},
	}
	for _, c := range cases {
		got, err := repo.IDByRef(c.ws, c.ref)
		if c.want == 0 {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("IDByRef(%d, %q) = %d, %v; want not found", c.ws, c.ref, got, err)
			}
			continue
		}
		if err != nil || got != c.want {
			t.Errorf("IDByRef(%d, %q) = %d, %v; want %d", c.ws, c.ref, got, err, c.want)
		}
	}
}
