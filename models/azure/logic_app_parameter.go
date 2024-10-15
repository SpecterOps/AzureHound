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

import (
	"github.com/certmichelin/azurehound/v2/enums"
)

type LogicAppParameter struct {
	Description string `json:"description,omitempty"`
	//Metadata - marked as object in MSDN, however no other description available - in testing was not able to return a value here
	Metadata interface{}         `json:"metadata,omitempty"`
	Type     enums.ParameterType `json:"type,omitempty"`
	Value    interface{}         `json:"value,omitempty"`
}

func (s LogicAppParameter) GetValue() interface{} {
	switch s.Type {
	case enums.ArrayType:
		return s.Value.([]interface{})
	case enums.BoolType:
		return s.Value.(bool)
	case enums.FloatType:
		return s.Value.(float64)
	case enums.IntType:
		return s.Value.(int)
	case enums.NotSpecifiedType:
		return s.Value
	case enums.ObjectType:
		return s.Value
	case enums.SecureObjectType:
		return s.Value
	case enums.SecureStringType:
		return s.Value
	case enums.StringType:
		return s.Value.(string)
	default:
		return s.Value
	}
}
