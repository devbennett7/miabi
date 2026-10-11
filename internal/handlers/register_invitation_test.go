// SPDX-FileCopyrightText: 2026 Jonas Kaninda
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/registration"
	"github.com/miabi-io/miabi/internal/services/settings"
	"github.com/miabi-io/miabi/internal/services/workspace"
	"github.com/miabi-io/miabi/internal/storage/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func invitationTestHandler(t *testing.T, signupEnabled string) *RegisterHandler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.WorkspaceInvitation{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	users := repositories.NewUserRepository(db)
	provider := settings.NewFixedProvider(map[string]models.Setting{
		settings.KeyInvitationSignupEnabled: {Key: settings.KeyInvitationSignupEnabled, Value: signupEnabled, Type: models.SettingTypeBool},
	})
	// Sign-up is closed: invitation sign-up must not depend on it.
	h := NewRegisterHandler(registration.NewService(false, provider, nil), nil, users, nil, nil)
	h.SetInvitations(workspace.NewService(repositories.NewWorkspaceRepository(db), users, nil), nil)
	return h
}

func TestRegisterFromInvitationRespectsTheSetting(t *testing.T) {
	h := invitationTestHandler(t, "false")
	rec := httptest.NewRecorder()
	c := okapi.NewContext(okapi.New(), rec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register/invitation", nil))
	req := &RegisterFromInvitationRequest{}
	req.Body.Token, req.Body.Name, req.Body.Password = "whatever", "Ada", "correct-horse-battery"
	if err := h.RegisterFromInvitation(c, req); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("status %d, want 403 while invitation sign-up is off", rec.Code)
	}
}

func TestRegisterFromInvitationRefusesAnUnknownToken(t *testing.T) {
	h := invitationTestHandler(t, "true")
	rec := httptest.NewRecorder()
	c := okapi.NewContext(okapi.New(), rec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register/invitation", nil))
	req := &RegisterFromInvitationRequest{}
	req.Body.Token, req.Body.Name, req.Body.Password = "not-a-token", "Ada", "correct-horse-battery"
	if err := h.RegisterFromInvitation(c, req); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400 for an unknown token", rec.Code)
	}
	if n, _ := h.users.Count(); n != 0 {
		t.Errorf("user count = %d, want 0", n)
	}
}

func TestPreviewInvitationAnswersOneWayForUnknownTokens(t *testing.T) {
	h := invitationTestHandler(t, "true")
	rec := httptest.NewRecorder()
	c := okapi.NewContext(okapi.New(), rec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/invitations/preview", nil))
	req := &InvitationTokenRequest{}
	req.Body.Token = "not-a-token"
	if err := h.PreviewInvitation(c, req); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("status %d, want 404", rec.Code)
	}
}

func TestSelfSignUpRecordsItsSource(t *testing.T) {
	h, users := registerTestHandler(t, "")
	if r := doRegister(t, h, "Ada", "ada@example.com", "correct-horse-battery"); r.status != http.StatusOK {
		t.Fatalf("sign-up: %d %s", r.status, r.body)
	}
	u, err := users.FindByEmail("ada@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.SignupSource != models.SignupSourceSelf {
		t.Errorf("signup source %q, want %q", u.SignupSource, models.SignupSourceSelf)
	}
}
