// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package update

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var dereferencableKindsNamedEntityConfig = map[reflect.Kind]struct{}{
	reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElementNamedEntityConfig(t reflect.Kind) bool {
	_, exists := dereferencableKindsNamedEntityConfig[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshalerHookNamedEntityConfig(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElementNamedEntityConfig(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

		raw, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		res := reflect.New(to).Interface()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			fmt.Printf("Failed to umarshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		return res, nil
	}

	return data, nil
}

func decode_NamedEntityConfig(input, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			jsonUnmarshalerHookNamedEntityConfig,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func join_NamedEntityConfig(arr interface{}, sep string) string {
	listValue := reflect.ValueOf(arr)
	strs := make([]string, 0, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		strs = append(strs, fmt.Sprintf("%v", listValue.Index(i)))
	}

	return strings.Join(strs, sep)
}

func testDecodeJson_NamedEntityConfig(t *testing.T, val, result interface{}) {
	assert.NoError(t, decode_NamedEntityConfig(val, result))
}

func testDecodeSlice_NamedEntityConfig(t *testing.T, vStringSlice, result interface{}) {
	assert.NoError(t, decode_NamedEntityConfig(vStringSlice, result))
}

func TestNamedEntityConfig_GetPFlagSet(t *testing.T) {
	val := NamedEntityConfig{}
	cmdFlags := val.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())
}

func TestNamedEntityConfig_SetFlags(t *testing.T) {
	actual := NamedEntityConfig{}
	cmdFlags := actual.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())

	t.Run("Test_archive", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vBool, err := cmdFlags.GetBool("archive"); err == nil {
				assert.Equal(t, bool(namedEntityConfig.Archive), vBool)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("archive", testValue)
			if vBool, err := cmdFlags.GetBool("archive"); err == nil {
				testDecodeJson_NamedEntityConfig(t, fmt.Sprintf("%v", vBool), &actual.Archive)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_activate", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vBool, err := cmdFlags.GetBool("activate"); err == nil {
				assert.Equal(t, bool(namedEntityConfig.Activate), vBool)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("activate", testValue)
			if vBool, err := cmdFlags.GetBool("activate"); err == nil {
				testDecodeJson_NamedEntityConfig(t, fmt.Sprintf("%v", vBool), &actual.Activate)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_description", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("description"); err == nil {
				assert.Equal(t, string(namedEntityConfig.Description), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("description", testValue)
			if vString, err := cmdFlags.GetString("description"); err == nil {
				testDecodeJson_NamedEntityConfig(t, fmt.Sprintf("%v", vString), &actual.Description)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
}