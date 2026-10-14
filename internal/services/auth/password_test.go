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
package auth

import (
	"errors"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func seedUser(t *testing.T, s *Service, username, password string) *models.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{Name: username, Username: username, Email: username + "@example.com", PasswordHash: string(hash), Role: models.SystemRoleUser, Active: true}
	if err := s.users.Create(&u); err != nil {
		t.Fatal(err)
	}
	return &u
}

// Sign-in keeps accepting passwords set under the earlier 8-character rule, and nothing shorter.
func TestLoginMinimumIsEight(t *testing.T) {
	s := authTestService(t)
	seedUser(t, s, "eight", "abcd1234")
	seedUser(t, s, "seven", "abc1234")

	if _, err := s.Authenticate("eight", "abcd1234"); err != nil {
		t.Fatalf("8-character password refused at sign-in: %v", err)
	}
	if _, err := s.Authenticate("seven", "abc1234"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("7-character password: err = %v, want ErrInvalidCredentials", err)
	}
}

func TestNewPasswordsNeedTwelve(t *testing.T) {
	s := authTestService(t)
	u := seedUser(t, s, "changer", "abcd1234")

	if err := s.ChangePassword(u.ID, "abcd1234", "elevenchars"); !errors.Is(err, models.ErrPasswordTooShort) {
		t.Fatalf("11-character new password: err = %v, want ErrPasswordTooShort", err)
	}
	if err := s.ChangePassword(u.ID, "abcd1234", "twelve-chars"); err != nil {
		t.Fatalf("12-character new password refused: %v", err)
	}
	if _, err := s.Authenticate("changer", "twelve-chars"); err != nil {
		t.Fatalf("sign-in with the new password: %v", err)
	}
	if _, err := s.ResetPassword("any-token", "short"); !errors.Is(err, models.ErrPasswordTooShort) {
		t.Fatalf("reset with a short password: err = %v, want ErrPasswordTooShort", err)
	}
}
