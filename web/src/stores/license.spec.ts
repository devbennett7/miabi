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

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const auth = { isAdmin: false }
const getLicense = vi.fn()
const entitlements = vi.fn()

vi.mock('@/stores/auth', () => ({ useAuthStore: () => auth }))
vi.mock('@/api/admin', () => ({ adminApi: { getLicense: () => getLicense() } }))
vi.mock('@/api/license', () => ({ licenseApi: { entitlements: () => entitlements() } }))

import { useLicenseStore } from './license'

const reply = <T>(data: T) => Promise.resolve({ data: { data } })

describe('license store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    getLicense.mockReset()
    entitlements.mockReset()
  })

  it('gives a non-admin the EE flags without calling the admin endpoint', async () => {
    auth.isAdmin = false
    entitlements.mockReturnValue(reply({ edition: 'enterprise', state: 'valid', flags: { advanced_canary: true } }))
    const store = useLicenseStore()
    await store.load()
    expect(getLicense).not.toHaveBeenCalled()
    expect(store.has('advanced_canary')).toBe(true)
    expect(store.mutable('advanced_canary')).toBe(true)
    expect(store.edition).toBe('enterprise')
  })

  it('keeps the full view for admins', async () => {
    auth.isAdmin = true
    getLicense.mockReturnValue(reply({ edition: 'enterprise', state: 'grace', flags: { custom_roles: true }, warnings: ['license_grace'] }))
    const store = useLicenseStore()
    await store.load()
    expect(entitlements).not.toHaveBeenCalled()
    expect(store.has('custom_roles')).toBe(true)
    expect(store.warnings).toEqual(['license_grace'])
  })

  it('reloads when the signed-in role changes in the same tab', async () => {
    auth.isAdmin = false
    entitlements.mockReturnValue(reply({ edition: 'enterprise', state: 'valid', flags: {} }))
    const store = useLicenseStore()
    await store.load()

    auth.isAdmin = true
    getLicense.mockReturnValue(reply({ edition: 'enterprise', state: 'valid', flags: {}, warnings: [] }))
    await store.load()
    expect(getLicense).toHaveBeenCalledTimes(1)
    expect(store.view).not.toBeNull()
  })
})
