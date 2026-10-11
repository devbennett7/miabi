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

package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/middlewares"
)

// fakeAppRefs resolves refs in workspace 3 only, so a lookup in another workspace misses.
type fakeAppRefs map[string]uint

func (f fakeAppRefs) IDByRef(workspaceID uint, ref string) (uint, error) {
	if id, ok := f[ref]; ok && workspaceID == 3 {
		return id, nil
	}
	return 0, errors.New("not found")
}

func TestAnalyticsAppFilterAcceptsIDUIDOrName(t *testing.T) {
	const uid = "01a12690-f759-7976-a4e7-d6b8a43884fa"
	h := &AnalyticsHandler{apps: fakeAppRefs{"gitea": 110, uid: 110}}
	app := okapi.New()
	app.Get("/analytics", func(c *okapi.Context) error {
		id, err := h.appFilterParam(c)
		if err != nil {
			return c.AbortNotFound("application not found")
		}
		if id == nil {
			return c.String(http.StatusOK, "all")
		}
		return c.String(http.StatusOK, strconv.FormatUint(uint64(*id), 10))
	}, okapi.UseMiddleware(func(c *okapi.Context) error {
		c.Set(middlewares.CtxWorkspaceID, 3)
		return c.Next()
	}))

	tests := []struct {
		query    string
		wantCode int
		wantBody string
	}{
		{"", http.StatusOK, "all"},
		{"?app=110", http.StatusOK, "110"},
		{"?app=999", http.StatusOK, "999"}, // a deleted app's traffic stays filterable by id
		{"?app=gitea", http.StatusOK, "110"},
		{"?app=" + uid, http.StatusOK, "110"},
		{"?app=missing", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/analytics"+tt.query, nil))
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d (%s)", rec.Code, tt.wantCode, rec.Body.String())
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
