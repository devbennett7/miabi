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

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type appComputeRow struct {
	ID          uint `gorm:"primaryKey"`
	WorkspaceID uint
	NanoCPUs    int64
	MemoryBytes int64
	Replicas    int
	RuntimeKind string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (appComputeRow) TableName() string { return "applications" }

func TestSumResourcesCountsServiceReplicas(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&appComputeRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	rows := []appComputeRow{
		{ID: 1, WorkspaceID: 1, NanoCPUs: 5e8, MemoryBytes: 512, Replicas: 4, RuntimeKind: "service"},
		{ID: 2, WorkspaceID: 1, NanoCPUs: 1e9, MemoryBytes: 256, Replicas: 3, RuntimeKind: "container"},
		{ID: 3, WorkspaceID: 1, NanoCPUs: 1e9, MemoryBytes: 100, Replicas: 0, RuntimeKind: "service"},
		{ID: 4, WorkspaceID: 2, NanoCPUs: 9e9, MemoryBytes: 999, Replicas: 1, RuntimeKind: "service"},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	repo := NewApplicationRepository(db)

	cpu, mem, err := repo.SumResourcesByWorkspace(1, 0)
	if err != nil {
		t.Fatalf("sum: %v", err)
	}
	if cpu != 4e9 || mem != 4*512+256+100 {
		t.Errorf("sum = %d nano / %d bytes, want 4e9 / %d", cpu, mem, 4*512+256+100)
	}

	cpu, mem, err = repo.SumResourcesByWorkspace(1, 1)
	if err != nil {
		t.Fatalf("sum excluding: %v", err)
	}
	if cpu != 2e9 || mem != 356 {
		t.Errorf("sum excluding the service = %d nano / %d bytes, want 2e9 / 356", cpu, mem)
	}
}
