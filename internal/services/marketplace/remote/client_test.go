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

package remote

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCapTransportFailsClosed(t *testing.T) {
	const max = 1024
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(bytes.Repeat([]byte("x"), max*4)) // four times the budget
	}))
	defer srv.Close()

	hc := &http.Client{Transport: &capTransport{base: http.DefaultTransport, max: max}}
	resp, err := hc.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	_, err = io.ReadAll(resp.Body)
	if !errors.Is(err, errBundleTooLarge) {
		t.Fatalf("expected errBundleTooLarge, got %v", err)
	}
}

func TestCapTransportPassesUnderLimit(t *testing.T) {
	const max = 1024
	body := bytes.Repeat([]byte("y"), max/2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	hc := &http.Client{Transport: &capTransport{base: http.DefaultTransport, max: max}}
	resp, err := hc.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read under limit should succeed: %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Fatal("body corrupted under the limit")
	}
}
