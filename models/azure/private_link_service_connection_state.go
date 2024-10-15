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

package azure

import "github.com/certmichelin/azurehound/v2/enums"

// An object that represents the approval state of the private link connection.
type PrivateLinkServiceConnectionState struct {
	// A message indicating if changes on the service provider require any updates on the consumer.
	ActionsRequired string `json:"actionsRequired,omitempty"`

	// The reason for approval or rejection.
	Description string `json:"description,omitempty"`

	// Indicates whether the connection has been approved, rejected or removed by the key vault owner.
	Status enums.EndpointConnectionStatus `json:"status,omitempty"`
}
