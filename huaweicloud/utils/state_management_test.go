package utils

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestStateManagementFunc_GetNestedObjectFromRawConfig(t *testing.T) {
	// Create a comprehensive test object with various data types and nesting levels
	complexObj := cty.ObjectVal(map[string]cty.Value{
		"string_field": cty.StringVal("simple_string"),
		"number_field": cty.NumberIntVal(42),
		"bool_field":   cty.BoolVal(true),
		"null_field":   cty.NullVal(cty.String),
		"list_field": cty.ListVal([]cty.Value{
			cty.StringVal("item1"),
			cty.StringVal("item2"),
			cty.StringVal("item3"),
		}),
		"nested_object": cty.ObjectVal(map[string]cty.Value{
			"level1_string": cty.StringVal("level1_value"),
			"level1_number": cty.NumberIntVal(100),
			"level1_list": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"id":   cty.NumberIntVal(1),
					"name": cty.StringVal("nested_item_1"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"id":   cty.NumberIntVal(2),
					"name": cty.StringVal("nested_item_2"),
				}),
			}),
		}),
		"deep_nested": cty.ObjectVal(map[string]cty.Value{
			"level1": cty.ObjectVal(map[string]cty.Value{
				"level2": cty.ObjectVal(map[string]cty.Value{
					"level3": cty.ObjectVal(map[string]cty.Value{
						"final_value": cty.StringVal("deep_nested_value"),
					}),
				}),
			}),
		}),
		"mixed_types": cty.ObjectVal(map[string]cty.Value{
			"strings": cty.ListVal([]cty.Value{
				cty.StringVal("str1"),
				cty.StringVal("str2"),
			}),
			"numbers": cty.ListVal([]cty.Value{
				cty.NumberIntVal(10),
				cty.NumberIntVal(20),
				cty.NumberIntVal(30),
			}),
			"booleans": cty.ListVal([]cty.Value{
				cty.BoolVal(true),
				cty.BoolVal(false),
			}),
		}),
	})

	tests := []struct {
		name     string
		input    cty.Value
		path     string
		expected interface{}
	}{
		// Test empty path - should return entire object
		{
			name:  "empty path returns entire object",
			input: complexObj,
			path:  "",
			expected: map[string]interface{}{
				"string_field": "simple_string",
				"number_field": float64(42),
				"bool_field":   true,
				"null_field":   nil,
				"list_field": []interface{}{
					"item1",
					"item2",
					"item3",
				},
				"nested_object": map[string]interface{}{
					"level1_string": "level1_value",
					"level1_number": float64(100),
					"level1_list": []interface{}{
						map[string]interface{}{
							"id":   float64(1),
							"name": "nested_item_1",
						},
						map[string]interface{}{
							"id":   float64(2),
							"name": "nested_item_2",
						},
					},
				},
				"deep_nested": map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": map[string]interface{}{
							"level3": map[string]interface{}{
								"final_value": "deep_nested_value",
							},
						},
					},
				},
				"mixed_types": map[string]interface{}{
					"strings": []interface{}{
						"str1",
						"str2",
					},
					"numbers": []interface{}{
						float64(10),
						float64(20),
						float64(30),
					},
					"booleans": []interface{}{
						true,
						false,
					},
				},
			},
		},
		// Test simple field access
		{
			name:     "simple string field access",
			input:    complexObj,
			path:     "string_field",
			expected: "simple_string",
		},
		{
			name:     "simple number field access",
			input:    complexObj,
			path:     "number_field",
			expected: float64(42),
		},
		{
			name:     "simple boolean field access",
			input:    complexObj,
			path:     "bool_field",
			expected: true,
		},
		{
			name:     "null field access",
			input:    complexObj,
			path:     "null_field",
			expected: nil,
		},
		// Test list access
		{
			name:     "list field access",
			input:    complexObj,
			path:     "list_field",
			expected: []interface{}{"item1", "item2", "item3"},
		},
		{
			name:     "list element access by index",
			input:    complexObj,
			path:     "list_field.0",
			expected: "item1",
		},
		{
			name:     "list element access by index 1",
			input:    complexObj,
			path:     "list_field.1",
			expected: "item2",
		},
		{
			name:     "list element access by index 2",
			input:    complexObj,
			path:     "list_field.2",
			expected: "item3",
		},
		// Test nested object access
		{
			name:  "nested object access",
			input: complexObj,
			path:  "nested_object",
			expected: map[string]interface{}{
				"level1_string": "level1_value",
				"level1_number": float64(100),
				"level1_list": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "nested_item_1",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "nested_item_2",
					},
				},
			},
		},
		{
			name:     "nested object property access",
			input:    complexObj,
			path:     "nested_object.level1_string",
			expected: "level1_value",
		},
		{
			name:     "nested object number property access",
			input:    complexObj,
			path:     "nested_object.level1_number",
			expected: float64(100),
		},
		// Test nested list access
		{
			name:  "nested list access",
			input: complexObj,
			path:  "nested_object.level1_list",
			expected: []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "nested_item_1",
				},
				map[string]interface{}{
					"id":   float64(2),
					"name": "nested_item_2",
				},
			},
		},
		{
			name:  "nested list element access",
			input: complexObj,
			path:  "nested_object.level1_list.0",
			expected: map[string]interface{}{
				"id":   float64(1),
				"name": "nested_item_1",
			},
		},
		{
			name:     "nested list element property access",
			input:    complexObj,
			path:     "nested_object.level1_list.0.name",
			expected: "nested_item_1",
		},
		{
			name:     "nested list element property access 2",
			input:    complexObj,
			path:     "nested_object.level1_list.1.id",
			expected: float64(2),
		},
		// Test deep nesting
		{
			name:  "deep nested access level1",
			input: complexObj,
			path:  "deep_nested.level1",
			expected: map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": map[string]interface{}{
						"final_value": "deep_nested_value",
					},
				},
			},
		},
		{
			name:  "deep nested access level2",
			input: complexObj,
			path:  "deep_nested.level1.level2",
			expected: map[string]interface{}{
				"level3": map[string]interface{}{
					"final_value": "deep_nested_value",
				},
			},
		},
		{
			name:  "deep nested access level3",
			input: complexObj,
			path:  "deep_nested.level1.level2.level3",
			expected: map[string]interface{}{
				"final_value": "deep_nested_value",
			},
		},
		{
			name:     "deep nested final value access",
			input:    complexObj,
			path:     "deep_nested.level1.level2.level3.final_value",
			expected: "deep_nested_value",
		},
		// Test mixed types
		{
			name:     "mixed types strings access",
			input:    complexObj,
			path:     "mixed_types.strings",
			expected: []interface{}{"str1", "str2"},
		},
		{
			name:     "mixed types numbers access",
			input:    complexObj,
			path:     "mixed_types.numbers",
			expected: []interface{}{float64(10), float64(20), float64(30)},
		},
		{
			name:     "mixed types booleans access",
			input:    complexObj,
			path:     "mixed_types.booleans",
			expected: []interface{}{true, false},
		},
		{
			name:     "mixed types specific string element",
			input:    complexObj,
			path:     "mixed_types.strings.0",
			expected: "str1",
		},
		{
			name:     "mixed types specific number element",
			input:    complexObj,
			path:     "mixed_types.numbers.1",
			expected: float64(20),
		},
		{
			name:     "mixed types specific boolean element",
			input:    complexObj,
			path:     "mixed_types.booleans.1",
			expected: false,
		},
		// Test edge cases and error conditions
		{
			name:     "non-existent path returns nil",
			input:    complexObj,
			path:     "non.existent.path",
			expected: nil,
		},
		{
			name:     "invalid list index returns nil",
			input:    complexObj,
			path:     "list_field.999",
			expected: nil,
		},
		{
			name:     "negative list index returns nil",
			input:    complexObj,
			path:     "list_field.-1",
			expected: nil,
		},
		{
			name:     "non-numeric list index returns nil",
			input:    complexObj,
			path:     "list_field.invalid",
			expected: nil,
		},
		{
			name:     "accessing property on primitive type returns primitive value",
			input:    complexObj,
			path:     "string_field.property",
			expected: "simple_string",
		},
		{
			name:     "accessing property on null value returns nil",
			input:    complexObj,
			path:     "null_field.property",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNestedObjectFromRawConfig(tt.input, tt.path)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetNestedObjectFromRawConfig(%v, %q) = %v, want %v",
					tt.input, tt.path, result, tt.expected)
			}
		})
	}
}

func TestStateManagementFunc_RefreshObjectParamOriginValues(t *testing.T) {
	// Test the core logic of RefreshObjectParamOriginValues by testing its components

	// Test 1: Test deepCopyInterface functionality
	originalData := map[string]interface{}{
		"policy": map[string]interface{}{
			"annotations": []interface{}{
				map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				map[string]interface{}{
					"key1": "value3",
					"key2": "value4",
				},
			},
			"config": map[string]interface{}{
				"setting1": "setting_value1",
				"setting2": "setting_value2",
			},
		},
	}

	// Test deep copy
	copiedData := deepCopyInterface(originalData)
	if !reflect.DeepEqual(originalData, copiedData) {
		t.Fatalf("deepCopyInterface failed to create identical copy")
	}

	// Modify copied data to verify deep copy
	copiedMap := copiedData.(map[string]interface{})
	policyMap := copiedMap["policy"].(map[string]interface{})
	annotations := policyMap["annotations"].([]interface{})
	firstAnnotation := annotations[0].(map[string]interface{})
	firstAnnotation["key1"] = "modified_value"

	// Verify original data was not affected
	originalPolicyMap := originalData["policy"].(map[string]interface{})
	originalAnnotations := originalPolicyMap["annotations"].([]interface{})
	originalFirstAnnotation := originalAnnotations[0].(map[string]interface{})

	if originalFirstAnnotation["key1"] != "value1" {
		t.Errorf("deepCopyInterface did not create a deep copy, modifying result affected input")
	}

	// Test 2: Test updateNestedStructureSafely functionality
	testCases := []struct {
		name        string
		current     interface{}
		parts       []string
		value       interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "update nested map key",
			current: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"key": "old_value",
					},
				},
			},
			parts: []string{"level1", "level2", "key"},
			value: "new_value",
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"key": "new_value",
					},
				},
			},
			expectError: false,
		},
		{
			name: "update list element",
			current: []interface{}{
				map[string]interface{}{
					"id":   1,
					"name": "old_name",
				},
				map[string]interface{}{
					"id":   2,
					"name": "item2",
				},
			},
			parts: []string{"0", "name"},
			value: "new_name",
			expected: []interface{}{
				map[string]interface{}{
					"id":   1,
					"name": "new_name",
				},
				map[string]interface{}{
					"id":   2,
					"name": "item2",
				},
			},
			expectError: false,
		},
		{
			name: "error on non-existent map key",
			current: map[string]interface{}{
				"existing_key": "value",
			},
			parts:       []string{"non_existent", "sub_key"},
			value:       "new_value",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a deep copy of current to avoid modifying the test data
			currentCopy := deepCopyInterface(tt.current)

			result, err := updateNestedStructureSafely(currentCopy, tt.parts, tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("updateNestedStructureSafely(%v, %v, %v) expected error but got none",
						tt.current, tt.parts, tt.value)
				}
				return
			}

			if err != nil {
				t.Errorf("updateNestedStructureSafely(%v, %v, %v) unexpected error: %v",
					tt.current, tt.parts, tt.value, err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("updateNestedStructureSafely(%v, %v, %v) = %v, want %v",
					tt.current, tt.parts, tt.value, result, tt.expected)
			}

			// Verify that the original current was not modified
			if !reflect.DeepEqual(tt.current, currentCopy) {
				t.Errorf("updateNestedStructureSafely modified the original input")
			}
		})
	}

	// Test 3: Test GetNestedObjectFromRawConfig with cty values
	// This simulates what RefreshObjectParamOriginValues would do
	ctyObj := cty.ObjectVal(map[string]cty.Value{
		"policy": cty.ObjectVal(map[string]cty.Value{
			"annotations": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"key1": cty.StringVal("value1"),
					"key2": cty.StringVal("value2"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"key1": cty.StringVal("value3"),
					"key2": cty.StringVal("value4"),
				}),
			}),
			"config": cty.ObjectVal(map[string]cty.Value{
				"setting1": cty.StringVal("setting_value1"),
				"setting2": cty.StringVal("setting_value2"),
			}),
		}),
	})

	// Test extracting nested objects
	annotationsResult := GetNestedObjectFromRawConfig(ctyObj, "policy.annotations")
	expectedAnnotations := []interface{}{
		map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
		map[string]interface{}{
			"key1": "value3",
			"key2": "value4",
		},
	}

	if !reflect.DeepEqual(annotationsResult, expectedAnnotations) {
		t.Errorf("GetNestedObjectFromRawConfig(policy.annotations) = %v, want %v",
			annotationsResult, expectedAnnotations)
	}

	configResult := GetNestedObjectFromRawConfig(ctyObj, "policy.config")
	expectedConfig := map[string]interface{}{
		"setting1": "setting_value1",
		"setting2": "setting_value2",
	}

	if !reflect.DeepEqual(configResult, expectedConfig) {
		t.Errorf("GetNestedObjectFromRawConfig(policy.config) = %v, want %v",
			configResult, expectedConfig)
	}

	// Test extracting non-existent path
	nonExistentResult := GetNestedObjectFromRawConfig(ctyObj, "policy.non_existent")
	if nonExistentResult != nil {
		t.Errorf("GetNestedObjectFromRawConfig(policy.non_existent) = %v, want nil",
			nonExistentResult)
	}

	t.Logf("All RefreshObjectParamOriginValues component tests passed successfully")
}
