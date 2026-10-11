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

package dto

// Response is the standard API response envelope with a generic data field.
type Response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

// PageableResponse is the paginated API response envelope.
type PageableResponse[T any] struct {
	Success  bool     `json:"success"`
	Data     []T      `json:"data"`
	Pageable Pageable `json:"pageable"`
}

// Pageable holds pagination metadata.
type Pageable struct {
	CurrentPage   int   `json:"current_page"`
	Size          int   `json:"size"`
	TotalPages    int   `json:"total_pages"`
	TotalElements int64 `json:"total_elements"`
	Empty         bool  `json:"empty"`
}

// ErrorResponseBody is the error envelope returned by the custom error handler.
type ErrorResponseBody struct {
	Success bool       `json:"success"`
	Data    any        `json:"data"`
	Error   *ErrorInfo `json:"error"`
}

// ErrorInfo holds error details.
type ErrorInfo struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// MessageData is a simple message payload for Response[MessageData].
type MessageData struct {
	Message string `json:"message"`
}
