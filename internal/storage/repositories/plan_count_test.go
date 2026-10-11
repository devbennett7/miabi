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
	"testing"

	"github.com/miabi-io/miabi/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newPlanRepo(t *testing.T) *PlanRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Plan{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	rows := []models.Plan{
		{Name: "Pro", IsActive: true, IsDefault: true},
		{Name: "Starter", IsActive: true},
		{Name: models.UnlimitedPlanName, IsActive: true, System: true},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	return NewPlanRepository(db)
}

func TestCountIncludesSystemPlans(t *testing.T) {
	repo := newPlanRepo(t)

	n, err := repo.Count()
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("Count() = %d, want 3 — the system plan counts toward the cap", n)
	}
}
