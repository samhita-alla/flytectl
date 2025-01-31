// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package get

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var dereferencableKindsLaunchPlanConfig = map[reflect.Kind]struct{}{
	reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElementLaunchPlanConfig(t reflect.Kind) bool {
	_, exists := dereferencableKindsLaunchPlanConfig[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshalerHookLaunchPlanConfig(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElementLaunchPlanConfig(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

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

func decode_LaunchPlanConfig(input, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			jsonUnmarshalerHookLaunchPlanConfig,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func join_LaunchPlanConfig(arr interface{}, sep string) string {
	listValue := reflect.ValueOf(arr)
	strs := make([]string, 0, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		strs = append(strs, fmt.Sprintf("%v", listValue.Index(i)))
	}

	return strings.Join(strs, sep)
}

func testDecodeJson_LaunchPlanConfig(t *testing.T, val, result interface{}) {
	assert.NoError(t, decode_LaunchPlanConfig(val, result))
}

func testDecodeSlice_LaunchPlanConfig(t *testing.T, vStringSlice, result interface{}) {
	assert.NoError(t, decode_LaunchPlanConfig(vStringSlice, result))
}

func TestLaunchPlanConfig_GetPFlagSet(t *testing.T) {
	val := LaunchPlanConfig{}
	cmdFlags := val.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())
}

func TestLaunchPlanConfig_SetFlags(t *testing.T) {
	actual := LaunchPlanConfig{}
	cmdFlags := actual.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())

	t.Run("Test_execFile", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("execFile"); err == nil {
				assert.Equal(t, string(launchPlanConfig.ExecFile), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("execFile", testValue)
			if vString, err := cmdFlags.GetString("execFile"); err == nil {
				testDecodeJson_LaunchPlanConfig(t, fmt.Sprintf("%v", vString), &actual.ExecFile)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_version", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("version"); err == nil {
				assert.Equal(t, string(launchPlanConfig.Version), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("version", testValue)
			if vString, err := cmdFlags.GetString("version"); err == nil {
				testDecodeJson_LaunchPlanConfig(t, fmt.Sprintf("%v", vString), &actual.Version)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_latest", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vBool, err := cmdFlags.GetBool("latest"); err == nil {
				assert.Equal(t, bool(launchPlanConfig.Latest), vBool)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("latest", testValue)
			if vBool, err := cmdFlags.GetBool("latest"); err == nil {
				testDecodeJson_LaunchPlanConfig(t, fmt.Sprintf("%v", vBool), &actual.Latest)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
}
