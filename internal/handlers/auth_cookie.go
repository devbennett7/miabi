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
	"net/http"
	"strings"
	"time"

	"github.com/jkaninda/okapi"
	"github.com/miabi-io/miabi/internal/middlewares"
	"github.com/miabi-io/miabi/internal/services/auth"
)

// setSessionCookie stores the JWT in an HttpOnly, SameSite=Strict cookie so a browser never
// exposes it to JavaScript. CLI and API clients keep using the token from the response body.
// Secure is set whenever the request is HTTPS, so the cookie still works over plain HTTP in dev.
func setSessionCookie(c *okapi.Context, token string) {
	http.SetCookie(c.ResponseWriter(), &http.Cookie{
		Name:     middlewares.SessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(auth.TokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   requestIsHTTPS(c),
		SameSite: http.SameSiteStrictMode,
	})
}

func clearSessionCookie(c *okapi.Context) {
	http.SetCookie(c.ResponseWriter(), &http.Cookie{
		Name:     middlewares.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   requestIsHTTPS(c),
		SameSite: http.SameSiteStrictMode,
	})
}

// oauthStateCookie binds an SSO authorization request to the browser that started it, so a callback
// carrying someone else's code and state (login-CSRF) is refused.
const oauthStateCookie = "miabi_oauth_state"

// oauthTwoFactorCookie carries a parked SSO sign-in to the second-factor step without exposing its id
// in a URL.
const oauthTwoFactorCookie = "miabi_oauth_2fa"

// setFlowCookie stores a short-lived HttpOnly cookie for the SSO flow. The state cookie must be Lax: it
// is read on the top-level redirect back from the identity provider, which Strict would withhold.
func setFlowCookie(c *okapi.Context, name, value string, ttl time.Duration, sameSite http.SameSite) {
	maxAge := int(ttl.Seconds())
	if value == "" {
		maxAge = -1
	}
	http.SetCookie(c.ResponseWriter(), &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   requestIsHTTPS(c),
		SameSite: sameSite,
	})
}

func cookieValue(c *okapi.Context, name string) string {
	ck, err := c.Request().Cookie(name)
	if err != nil {
		return ""
	}
	return ck.Value
}

func requestIsHTTPS(c *okapi.Context) bool {
	if c.Request().TLS != nil {
		return true
	}
	return strings.EqualFold(c.Header("X-Forwarded-Proto"), "https")
}
