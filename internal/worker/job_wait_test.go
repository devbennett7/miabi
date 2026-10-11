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

package worker

import (
	"context"
	"testing"
	"time"

	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/storage/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newWaitHandler(t *testing.T) (*JobHandler, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.Exec(`CREATE TABLE deployments (id INTEGER PRIMARY KEY, application_id INTEGER, status TEXT, number INTEGER)`)
	db.Exec(`CREATE TABLE releases (id INTEGER PRIMARY KEY, application_id INTEGER, image TEXT, active BOOLEAN)`)
	if err := db.AutoMigrate(&models.Job{}); err != nil {
		t.Fatal(err)
	}
	h := &JobHandler{jobs: repositories.NewJobRepository(db)}
	h.SetDeployWait(repositories.NewDeploymentRepository(db), repositories.NewReleaseRepository(db))
	h.waitPoll = 10 * time.Millisecond
	return h, db
}

func waitingJob(t *testing.T, db *gorm.DB) *models.Job {
	t.Helper()
	dep := uint(9)
	j := &models.Job{WorkspaceID: 1, ApplicationID: 1, Command: []string{"x"}, Image: "old", Status: models.JobPending, WaitDeploymentID: &dep}
	if err := db.Create(j).Error; err != nil {
		t.Fatal(err)
	}
	return j
}

// The run starts once its deploy ends, on the image of the release it activated.
func TestAwaitDeployUsesNewRelease(t *testing.T) {
	h, db := newWaitHandler(t)
	db.Exec(`INSERT INTO deployments (id, application_id, status, number) VALUES (9, 1, 'deploying', 3)`)
	j := waitingJob(t, db)
	go func() {
		time.Sleep(30 * time.Millisecond)
		db.Exec(`INSERT INTO releases (application_id, image, active) VALUES (1, 'new', 1)`)
		db.Exec(`UPDATE deployments SET status = 'succeeded' WHERE id = 9`)
	}()
	if !h.awaitDeploy(context.Background(), j) || j.Image != "new" {
		t.Errorf("go on = false or image = %q, want the new release", j.Image)
	}
}

func TestAwaitDeploySkipsOnFailure(t *testing.T) {
	h, db := newWaitHandler(t)
	db.Exec(`INSERT INTO deployments (id, application_id, status, number) VALUES (9, 1, 'failed', 3)`)
	j := waitingJob(t, db)
	if h.awaitDeploy(context.Background(), j) {
		t.Fatal("a failed deploy should stop the run")
	}
	var got models.Job
	db.First(&got, j.ID)
	if got.Status != models.JobSkipped {
		t.Errorf("status = %s, want skipped", got.Status)
	}
}
