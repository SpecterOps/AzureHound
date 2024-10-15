package models_test

import (
	"encoding/json"
	"testing"

	"github.com/certmichelin/azurehound/v2/models"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	Bar string
}

func TestStripEmptyEntries(t *testing.T) {
	t.Run("should omit empty basic types", func(t *testing.T) {
		var pointer *Foo

		data := map[string]any{
			"array":   [0]any{},
			"map":     map[string]any{},
			"slice":   []any{},
			"string":  "",
			"bool":    false,
			"float32": float32(0),
			"float64": float64(0),
			"int":     int(0),
			"int8":    int8(0),
			"int16":   int16(0),
			"int32":   int32(0),
			"int64":   int64(0),
			"uint":    uint(0),
			"uint8":   uint8(0),
			"uint16":  uint16(0),
			"uint32":  uint32(0),
			"uint64":  uint64(0),
			"pointer": pointer,
			"nil":     nil,
		}

		require.NotEmpty(t, data)
		models.StripEmptyEntries(data)
		require.Empty(t, data)
	})

	t.Run("should not omit non-empty basic types", func(t *testing.T) {
		data := map[string]any{
			"array":   [1]any{1},
			"map":     map[any]any{"foo": "bar"},
			"slice":   []any{1},
			"string":  "foo",
			"bool":    true,
			"float32": float32(1),
			"float64": float64(1),
			"int":     int(1),
			"int8":    int8(1),
			"int16":   int16(1),
			"int32":   int32(1),
			"int64":   int64(1),
			"uint":    uint(1),
			"uint8":   uint8(1),
			"uint16":  uint16(1),
			"uint32":  uint32(1),
			"uint64":  uint64(1),
			"pointer": &Foo{},
		}

		require.NotEmpty(t, data)
		numKeys := len(data)
		models.StripEmptyEntries(data)
		require.NotEmpty(t, data)
		require.Equal(t, numKeys, len(data))
	})

	t.Run("should not omit empty struct types", func(t *testing.T) {
		data := map[string]any{
			"struct": Foo{Bar: "baz"},
		}

		require.NotEmpty(t, data)
		numKeys := len(data)
		models.StripEmptyEntries(data)
		require.NotEmpty(t, data)
		require.Equal(t, numKeys, len(data))
		require.NotEmpty(t, data["struct"])
		require.Equal(t, data["struct"].(Foo).Bar, "baz")
	})

	t.Run("should recursively strip non-empty, nested map[string]any entries", func(t *testing.T) {
		data := map[string]any{
			"empty": map[string]any{
				"false":       false,
				"emptystring": "",
				"emptynest": map[string]any{
					"false":       false,
					"emptystring": "",
				},
			},
			"nonempty": map[string]any{
				"emptyprop":    0,
				"nonemptyprop": 1,
			},
		}

		models.StripEmptyEntries(data)
		require.Nil(t, data["empty"])
		require.NotNil(t, data["nonempty"])
		require.IsType(t, map[string]any{}, data["nonempty"])
		nested := data["nonempty"].(map[string]any)
		require.Equal(t, 1, len(nested))
		require.Nil(t, nested["emptyprop"])
		require.Equal(t, 1, nested["nonemptyprop"])
	})

	t.Run("should strip empty slice entries of type map[string]any", func(t *testing.T) {
		data := map[string]any{
			"empty": []any{
				map[string]any{
					"false":       false,
					"emptystring": "",
				},
			},
			"emptynestedslice": []any{
				map[string]any{
					"nestedslice": []any{
						map[string]any{
							"false":       false,
							"emptystring": "",
						},
					},
					"emptystring": "",
				},
			},
			"nonempty": []any{
				map[string]any{
					"emptyprop":    0,
					"nonemptyprop": 1,
				},
			},
		}

		models.StripEmptyEntries(data)
		require.Nil(t, data["empty"])
		require.Nil(t, data["emptynestedslice"])
		require.NotNil(t, data["nonempty"])
		require.IsType(t, []any{}, data["nonempty"])
		slice := data["nonempty"].([]any)
		require.IsType(t, map[string]any{}, slice[0])
		entry := slice[0].(map[string]any)
		require.Nil(t, entry["emptyprop"])
		require.Equal(t, 1, entry["nonemptyprop"])
	})
}

func TestOmitEmpty(t *testing.T) {
	t.Run("should omit empty basic types", func(t *testing.T) {
		data := json.RawMessage(`{
			"string": "",
			"number": 0,
			"object": {},
			"array": [],
			"boolean": false,
			"null": null
		}`)

		filtered, err := models.OmitEmpty(data)
		require.Nil(t, err)
		require.Equal(t, `{}`, string(filtered))
	})

	t.Run("should not omit non-empty basic types except empty structs", func(t *testing.T) {
		data := json.RawMessage(`{
			"string": "foo",
			"number": 1,
			"object": { "bar": "" },
			"array": [1],
			"boolean": true
		}`)

		filtered, err := models.OmitEmpty(data)
		require.Nil(t, err)
		require.Equal(t, `{"array":[1],"boolean":true,"number":1,"string":"foo"}`, string(filtered))
	})

	t.Run("should omit empty struct/object types, just their empty properties", func(t *testing.T) {
		data := json.RawMessage(`{
			"object": { "bar": "" }
		}`)

		filtered, err := models.OmitEmpty(data)
		require.Nil(t, err)
		require.Equal(t, `{}`, string(filtered))
	})

	t.Run("should recursively strip non-empty, nested object entries", func(t *testing.T) {
		data := json.RawMessage(`{
			"empty": {},
			"nonempty": {
				"emptyprop": 0,
				"nonemptyprop": 1
			}
		}`)

		filtered, err := models.OmitEmpty(data)
		require.Nil(t, err)
		require.Equal(t, `{"nonempty":{"nonemptyprop":1}}`, string(filtered))
	})

	t.Run("should strip non-empty array entries of type object", func(t *testing.T) {
		data := json.RawMessage(`{
			"empty": [],
			"nonempty": [{
				"emptyprop": 0,
				"nonemptyprop": 1
			}]
		}`)

		filtered, err := models.OmitEmpty(data)
		require.Nil(t, err)
		require.Equal(t, `{"nonempty":[{"nonemptyprop":1}]}`, string(filtered))
	})

	t.Run("should return an error", func(t *testing.T) {
		invalidJson := json.RawMessage(`{]}`)
		_, err := models.OmitEmpty(invalidJson)
		require.Error(t, err)
	})
}
