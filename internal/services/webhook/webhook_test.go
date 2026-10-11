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

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/miabi-io/miabi/internal/models"
)

func TestSignIsVerifiableHMAC(t *testing.T) {
	secret := "topsecret"
	body := []byte(`{"event":"deploy.succeeded"}`)

	got := Sign(secret, body)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	want := hex.EncodeToString(mac.Sum(nil))
	if got != want {
		t.Fatalf("Sign() = %q, want %q", got, want)
	}
	if Sign("other", body) == got {
		t.Fatal("signature should differ with a different key")
	}
}

func TestBuildPayload(t *testing.T) {
	e := &models.AppEvent{
		WorkspaceID:   7,
		ApplicationID: 3,
		Type:          models.EventDeploySucceeded,
		Severity:      models.SeverityInfo,
		Message:       "deployed",
		Metadata:      map[string]string{"tag": "v1"},
		CreatedAt:     time.Unix(1700000000, 0).UTC(),
	}
	raw, err := BuildPayload(e)
	if err != nil {
		t.Fatalf("BuildPayload() error = %v", err)
	}
	var p Payload
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Event != string(models.EventDeploySucceeded) {
		t.Errorf("Event = %q, want %q", p.Event, models.EventDeploySucceeded)
	}
	if p.WorkspaceID != 7 || p.ApplicationID != 3 {
		t.Errorf("ids = (%d,%d), want (7,3)", p.WorkspaceID, p.ApplicationID)
	}
	if p.Metadata["tag"] != "v1" {
		t.Errorf("metadata not preserved: %v", p.Metadata)
	}
	if p.Timestamp == "" {
		t.Error("timestamp should be set")
	}
}

func TestValidateURL(t *testing.T) {
	cases := map[string]bool{
		"https://example.com/hook": true,
		"http://1.2.3.4:9000":      true,
		"":                         false,
		"ftp://example.com":        false,
		"not-a-url":                false,
	}
	for in, ok := range cases {
		err := validateURL(in)
		if ok && err != nil {
			t.Errorf("validateURL(%q) = %v, want nil", in, err)
		}
		if !ok && err == nil {
			t.Errorf("validateURL(%q) = nil, want error", in)
		}
	}
}

func TestValidateEvents(t *testing.T) {
	if err := validateEvents([]string{string(models.EventDeployFailed)}); err != nil {
		t.Errorf("notifiable event rejected: %v", err)
	}
	if err := validateEvents([]string{string(models.EventEnvUpdated)}); err == nil {
		t.Error("non-notifiable event should be rejected")
	}
}
