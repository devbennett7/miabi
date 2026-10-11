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

package notify

import (
	"strings"
	"testing"

	"github.com/miabi-io/miabi/internal/models"
)

func TestFormatMessage(t *testing.T) {
	e := &models.AppEvent{
		ApplicationID: 12,
		Type:          models.EventDeployFailed,
		Message:       "image pull failed",
	}
	msg := formatMessage(e)
	if !strings.Contains(msg, "Deployment failed") {
		t.Errorf("missing title in %q", msg)
	}
	if !strings.Contains(msg, "#12") {
		t.Errorf("missing application id in %q", msg)
	}
	if !strings.Contains(msg, "image pull failed") {
		t.Errorf("missing message in %q", msg)
	}
}

func TestValidateEvents(t *testing.T) {
	if err := validateEvents([]string{string(models.EventContainerDied)}); err != nil {
		t.Errorf("notifiable event rejected: %v", err)
	}
	if err := validateEvents([]string{"bogus.event"}); err == nil {
		t.Error("unknown event should be rejected")
	}
}

func TestRegistrySupportedChannels(t *testing.T) {
	r := NewRegistry()
	for _, typ := range []models.NotificationChannelType{
		models.ChannelTelegram, models.ChannelSlack, models.ChannelDiscord,
	} {
		if _, ok := r.Sender(typ); !ok {
			t.Errorf("%s sender not registered", typ)
		}
	}
	if _, ok := r.Sender(models.NotificationChannelType("teams")); ok {
		t.Fatal("unexpected sender for unsupported type")
	}
}
