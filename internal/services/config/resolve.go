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
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// refRe matches `${{ secrets.NAME }}` and `${{ env.NAME }}`, the spelling env
// values already use for secrets.
var refRe = regexp.MustCompile(`\$\{\{\s*(secrets|env)\.([A-Za-z0-9_.\-]+)\s*\}\}`)

// Resolver supplies the values a config's references stand for, scoped to the
// app the config is mounted into.
type Resolver interface {
	Secret(name string) (string, bool)
	Env(name string) (string, bool)
}

// ResolverFunc adapts two lookups into a Resolver.
type ResolverFunc struct {
	SecretFn func(string) (string, bool)
	EnvFn    func(string) (string, bool)
}

func (r ResolverFunc) Secret(name string) (string, bool) {
	if r.SecretFn == nil {
		return "", false
	}
	return r.SecretFn(name)
}

func (r ResolverFunc) Env(name string) (string, bool) {
	if r.EnvFn == nil {
		return "", false
	}
	return r.EnvFn(name)
}

// HasReferences reports whether content carries any reference.
func HasReferences(content string) bool { return refRe.MatchString(content) }

// SecretRefNames lists the secrets a config's files reference, so rotating one
// can find the apps mounting it.
func SecretRefNames(data map[string]string) []string {
	seen := map[string]bool{}
	for _, content := range data {
		for _, m := range refRe.FindAllStringSubmatch(content, -1) {
			if m[1] == "secrets" {
				seen[m[2]] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Resolve substitutes every reference in content. A secret of that name wins
// over an env var of the same name, whichever namespace was written, so a value
// promoted from env to the vault needs no edit here.
//
// An unresolvable reference fails: a file mounted with the placeholder still in
// it is worse than a deploy that stops.
func Resolve(content string, r Resolver) (string, error) {
	if r == nil || !refRe.MatchString(content) {
		return content, nil
	}

	var missing []string
	out := refRe.ReplaceAllStringFunc(content, func(match string) string {
		m := refRe.FindStringSubmatch(match)
		kind, name := m[1], m[2]

		if v, ok := r.Secret(name); ok {
			return v
		}
		if kind == "env" {
			if v, ok := r.Env(name); ok {
				return v
			}
		}
		missing = append(missing, kind+"."+name)
		return match
	})

	if len(missing) > 0 {
		return "", fmt.Errorf("%w: %s", ErrUnresolvedRef, strings.Join(dedupe(missing), ", "))
	}
	return out, nil
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
