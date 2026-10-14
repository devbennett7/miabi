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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/storage/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAdminCreateUserPassword(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatal(err)
	}
	users := repositories.NewUserRepository(db)
	app := okapi.New()
	app.Post("/admin/users", okapi.H(NewAdminUserHandler(db, users, nil, nil, nil, nil, nil, nil, 0).Create))

	post := func(body string) (int, map[string]any) {
		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		var out struct {
			Data map[string]any `json:"data"`
		}
		_ = json.Unmarshal(rec.Body.Bytes(), &out)
		return rec.Code, out.Data
	}

	code, data := post(`{"name":"Gen","email":"gen@example.com"}`)
	pw, _ := data["password"].(string)
	if code != http.StatusCreated || len(pw) != 20 {
		t.Fatalf("omitted password: status %d, password %q; want 201 and a generated 20-character password", code, pw)
	}
	u, err := users.FindByEmail("gen@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw)) != nil || !u.MustChangePassword {
		t.Errorf("generated password not stored as a temporary password: must_change=%v", u.MustChangePassword)
	}
	if _, leaked := data["password_hash"]; leaked {
		t.Error("response carries the password hash")
	}

	if code, data := post(`{"name":"Set","email":"set@example.com","password":"twelve-chars"}`); code != http.StatusCreated || data["password"] != nil {
		t.Errorf("chosen password: status %d, echoed password %v; want 201 and nothing echoed", code, data["password"])
	}
	if code, _ := post(`{"name":"Short","email":"short@example.com","password":"elevenchars"}`); code != http.StatusBadRequest {
		t.Errorf("11-character password: status %d, want 400", code)
	}
}
