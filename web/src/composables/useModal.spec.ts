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

import { describe, expect, it } from 'vitest'
import { modalKeyAction } from './useModal'

// A plain dialog: no region claims Tab or Escape.
const plain = { escapable: true, inEscapeThrough: false, inTabThrough: false }
// A shell terminal: it owns both keys.
const shell = { escapable: true, inEscapeThrough: true, inTabThrough: true }

describe('modalKeyAction', () => {
  it('closes on Escape in an ordinary dialog', () => {
    expect(modalKeyAction({ key: 'Escape', shiftKey: false, ...plain })).toBe('close')
  })

  it('leaves Escape alone while the dialog is mid-action', () => {
    expect(
      modalKeyAction({ key: 'Escape', shiftKey: false, ...plain, escapable: false }),
    ).toBe('ignore')
  })


  it('gives Escape to content that owns it', () => {
    expect(modalKeyAction({ key: 'Escape', shiftKey: false, ...shell })).toBe('ignore')
  })

  it('gives Tab to content that owns it, so a shell can complete a path', () => {
    expect(modalKeyAction({ key: 'Tab', shiftKey: false, ...shell })).toBe('ignore')
  })


  it('always runs the focus trap on Shift+Tab, even inside that content', () => {
    expect(modalKeyAction({ key: 'Tab', shiftKey: true, ...shell })).toBe('trap-tab')
  })

  it('traps Tab in an ordinary dialog', () => {
    expect(modalKeyAction({ key: 'Tab', shiftKey: false, ...plain })).toBe('trap-tab')
    expect(modalKeyAction({ key: 'Tab', shiftKey: true, ...plain })).toBe('trap-tab')
  })

  it('ignores every other key', () => {
    for (const key of ['a', 'Enter', 'ArrowUp', 'Control', ':', 'w']) {
      expect(modalKeyAction({ key, shiftKey: false, ...shell })).toBe('ignore')
    }
  })
})
