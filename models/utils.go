package models

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// upperStrings returns a new slice with each element uppercased. A nil input
// returns nil so empty slices remain omitted by omitempty json tags.
func upperStrings(values []string) []string {
	if values == nil {
		return nil
	}
	upper := make([]string, len(values))
	for i, value := range values {
		upper[i] = strings.ToUpper(value)
	}
	return upper
}

// UpperRoleAssignment returns a copy of the provided RoleAssignment with the
// Properties.PrincipalId and Properties.Scope uppercased. BloodHound ingest
// uppercases the principal for the edge endpoint and, for scope-matched
// convertors, compares the (now uppercased) scope against the likewise
// uppercased target id. The input is not mutated.
func UpperRoleAssignment(assignment azure.RoleAssignment) azure.RoleAssignment {
	assignment.Properties.PrincipalId = strings.ToUpper(assignment.Properties.PrincipalId)
	assignment.Properties.Scope = strings.ToUpper(assignment.Properties.Scope)
	return assignment
}

// UpperManagedIdentity returns a copy of the provided ManagedIdentity with the
// system-assigned PrincipalId and each user-assigned identity PrincipalId
// uppercased. The input is not mutated: the UserAssignedIdentities map is
// rebuilt into a fresh map.
func UpperManagedIdentity(identity azure.ManagedIdentity) azure.ManagedIdentity {
	identity.PrincipalId = strings.ToUpper(identity.PrincipalId)
	if identity.UserAssignedIdentities != nil {
		uais := make(map[string]azure.UserAssignedIdentity, len(identity.UserAssignedIdentities))
		for key, uai := range identity.UserAssignedIdentities {
			uai.PrincipalId = strings.ToUpper(uai.PrincipalId)
			uais[key] = uai
		}
		identity.UserAssignedIdentities = uais
	}
	return identity
}

func OmitEmpty(raw json.RawMessage) (json.RawMessage, error) {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	} else {
		StripEmptyEntries(data)
		return json.Marshal(data)
	}
}

// OmitEmptyUpper behaves like OmitEmpty but also uppercases the string values of
// the provided top-level keys. It is non-mutating with respect to the input:
// the raw message is unmarshaled into a fresh map that is edited and re-marshaled.
// Missing keys and non-string values are left untouched.
func OmitEmptyUpper(raw json.RawMessage, keys ...string) (json.RawMessage, error) {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	} else {
		StripEmptyEntries(data)
		for _, key := range keys {
			if value, ok := data[key].(string); ok {
				data[key] = strings.ToUpper(value)
			}
		}
		return json.Marshal(data)
	}
}

func StripEmptyEntries(data map[string]any) {
	for key, value := range data {
		if isEmpty(reflect.ValueOf(value)) {
			delete(data, key)
		} else if nested, ok := value.(map[string]any); ok { // recursively strip nested maps
			StripEmptyEntries(nested)
		} else if slice, ok := value.([]any); ok {
			value = make([]any, len(value.([]any)))
			i := 0
			for _, item := range slice {
				if mapValue, ok := item.(map[string]any); ok {
					StripEmptyEntries(mapValue)
				}
				if !isEmpty(reflect.ValueOf(item)) {
					value.([]any)[i] = item
					i++
				}
			}
			value = value.([]any)[:i]
		}

		// Strip top level if empty post recursive strip
		if _, ok := data[key]; ok && isEmpty(reflect.ValueOf(value)) {
			delete(data, key)
		}
	}
}

func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return value.Bool() == false
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Interface, reflect.Pointer:
		return value.IsNil()
	case reflect.Invalid:
		return true
	default:
		return false
	}
}
