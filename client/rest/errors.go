// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rest

import (
	"errors"
	"fmt"
)

// GraphError is a structured representation of a Microsoft Graph error
// response body ({"error": {"code": "...", "message": "..."}}). It is
// returned by the REST layer when a 4xx response decodes to that shape so
// callers can programmatically detect specific error codes.
type GraphError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *GraphError) Error() string {
	return fmt.Sprintf("graph error %d: %s - %s", e.StatusCode, e.Code, e.Message)
}

// IsExpiredPageToken reports whether err (or anything it wraps) is a Graph
// error with code "Directory_ExpiredPageToken", indicating the pagination
// cursor in @odata.nextLink has expired and the enumeration must be
// restarted from scratch.
func IsExpiredPageToken(err error) bool {
	var ge *GraphError
	return errors.As(err, &ge) && ge.Code == "Directory_ExpiredPageToken"
}
