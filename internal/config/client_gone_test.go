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

package config

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jkaninda/okapi"
)

func TestIgnoreClientGone(t *testing.T) {
	app := okapi.New()
	app.Use(ignoreClientGone)
	app.Get("/stream", func(c *okapi.Context) error { return context.Canceled })
	app.Get("/broken", func(c *okapi.Context) error { return errors.New("boom") })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/stream", nil).WithContext(ctx))
	if rec.Code == http.StatusInternalServerError {
		t.Errorf("a client that went away was answered as a server error")
	}

	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/broken", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("a real failure = %d, want 500", rec.Code)
	}
}
