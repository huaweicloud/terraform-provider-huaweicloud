package utils

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestStateManagementFunc_GetObjectFromRawConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected interface{}
	}{
		{
			name:     "null value",
			input:    cty.NullVal(cty.String),
			expected: nil,
		},
		{
			name:     "string value",
			input:    cty.StringVal("test"),
			expected: "test",
		},
		{
			name:     "number value",
			input:    cty.NumberIntVal(42),
			expected: float64(42),
		},
		{
			name:     "float value",
			input:    cty.NumberFloatVal(3.14),
			expected: 3.14,
		},
		{
			name:     "bool value",
			input:    cty.BoolVal(true),
			expected: true,
		},
		{
			name: "simple object",
			input: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
				"num": cty.NumberIntVal(123),
			}),
			expected: map[string]interface{}{
				"foo": "bar",
				"num": float64(123),
			},
		},
		{
			name: "nested object",
			input: cty.ObjectVal(map[string]cty.Value{
				"config": cty.ObjectVal(map[string]cty.Value{
					"key1": cty.StringVal("value1"),
					"key2": cty.StringVal("value2"),
				}),
			}),
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
		{
			name: "list of objects",
			input: cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("ZhangSan"),
					"age":  cty.NumberIntVal(18),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("LiSi"),
					"age":  cty.NumberIntVal(19),
				}),
			}),
			expected: []interface{}{
				map[string]interface{}{
					"name": "ZhangSan",
					"age":  float64(18),
				},
				map[string]interface{}{
					"name": "LiSi",
					"age":  float64(19),
				},
			},
		},
		{
			name: "complex nested structure",
			input: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
				"owners": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"name": cty.StringVal("ZhangSan"),
						"age":  cty.NumberIntVal(18),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"name": cty.StringVal("LiSi"),
						"age":  cty.NumberIntVal(19),
					}),
				}),
				"config": cty.ObjectVal(map[string]cty.Value{
					"key1": cty.StringVal("value1"),
					"key2": cty.StringVal("value2"),
				}),
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
			}),
			expected: map[string]interface{}{
				"foo": "bar",
				"owners": []interface{}{
					map[string]interface{}{
						"name": "ZhangSan",
						"age":  float64(18),
					},
					map[string]interface{}{
						"name": "LiSi",
						"age":  float64(19),
					},
				},
				"config": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
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
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetObjectFromRawConfig(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("[Check %d] The processing result of the GetObjectFromRawConfig method is not as expected, want %s, but got %s",
					i+1, Green(tt.expected), Yellow(result))
			}
		})
	}

	t.Logf("All processing results of the GetObjectFromRawConfig method meets expectation")
}

func TestStateManagementFunc_GetNestedObjectFromRawConfig(t *testing.T) {
	// 创建包含policy嵌套结构的对象
	complexObj := cty.ObjectVal(map[string]cty.Value{
		"foo": cty.StringVal("bar"),
		"owners": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("ZhangSan"),
				"age":  cty.NumberIntVal(18),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("LiSi"),
				"age":  cty.NumberIntVal(19),
			}),
		}),
		"config": cty.ObjectVal(map[string]cty.Value{
			"key1": cty.StringVal("value1"),
			"key2": cty.StringVal("value2"),
		}),
		"policy": cty.ObjectVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
			"owners": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("ZhangSan"),
					"age":  cty.NumberIntVal(18),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("LiSi"),
					"age":  cty.NumberIntVal(19),
				}),
			}),
			"config": cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
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
		}),
	})

	tests := []struct {
		name     string
		input    cty.Value
		path     string
		expected interface{}
	}{
		{
			name:     "simple property access",
			input:    complexObj,
			path:     "foo",
			expected: "bar",
		},
		{
			name:     "nested object access",
			input:    complexObj,
			path:     "config.key1",
			expected: "value1",
		},
		{
			name:     "list element access",
			input:    complexObj,
			path:     "owners.0.name",
			expected: "ZhangSan",
		},
		{
			name:     "list element access with number",
			input:    complexObj,
			path:     "owners.1.age",
			expected: float64(19),
		},
		{
			name:     "non-existent path",
			input:    complexObj,
			path:     "non.existent.path",
			expected: nil,
		},
		{
			name:  "empty path",
			input: complexObj,
			path:  "",
			expected: map[string]interface{}{
				"foo": "bar",
				"owners": []interface{}{
					map[string]interface{}{
						"name": "ZhangSan",
						"age":  float64(18),
					},
					map[string]interface{}{
						"name": "LiSi",
						"age":  float64(19),
					},
				},
				"config": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				"policy": map[string]interface{}{
					"foo": "bar",
					"owners": []interface{}{
						map[string]interface{}{
							"name": "ZhangSan",
							"age":  float64(18),
						},
						map[string]interface{}{
							"name": "LiSi",
							"age":  float64(19),
						},
					},
					"config": map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
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
				},
			},
		},
		{
			name:  "get policy object",
			input: complexObj,
			path:  "policy",
			expected: map[string]interface{}{
				"foo": "bar",
				"owners": []interface{}{
					map[string]interface{}{
						"name": "ZhangSan",
						"age":  float64(18),
					},
					map[string]interface{}{
						"name": "LiSi",
						"age":  float64(19),
					},
				},
				"config": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
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
			},
		},
		{
			name:  "get annotations list",
			input: complexObj,
			path:  "policy.annotations",
			expected: []interface{}{
				map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				map[string]interface{}{
					"key1": "value3",
					"key2": "value4",
				},
			},
		},
		{
			name:  "get first annotation",
			input: complexObj,
			path:  "policy.annotations.0",
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "get second annotation key1",
			input:    complexObj,
			path:     "policy.annotations.1.key1",
			expected: "value3",
		},
		{
			name:     "get config key1",
			input:    complexObj,
			path:     "policy.config.key1",
			expected: "value1",
		},
		{
			name:     "get first owner name",
			input:    complexObj,
			path:     "policy.owners.0.name",
			expected: "ZhangSan",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNestedObjectFromRawConfig(tt.input, tt.path)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("[Check %d] The processing result of the GetNestedObjectFromRawConfig method is not as expected, want %s, but got %s",
					i+1, Green(tt.expected), Yellow(result))
			}
		})
	}

	// 演示GetAttrFromRawConfig的限制
	t.Run("GetAttrFromRawConfig limitation", func(t *testing.T) {
		// 这只能获取到policy对象本身，不能获取annotations
		policyObj := complexObj.GetAttr("policy")
		if policyObj.IsNull() {
			t.Fatalf("[Check 1] The GetAttr() method should return the policy object, but got null")
		}

		// 验证policyObj是cty.Value类型，需要进一步处理
		policyMap := GetObjectFromRawConfig(policyObj)
		annotations := policyMap.(map[string]interface{})["annotations"]

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

		if !reflect.DeepEqual(annotations, expectedAnnotations) {
			t.Fatalf("[Check 2] The processing result of the annotations extraction is not as expected, want %s, but got %s",
				Green(expectedAnnotations), Yellow(annotations))
		}

		// 演示需要两步操作才能获取annotations
		annotationsObj := policyObj.GetAttr("annotations")
		annotationsList := GetObjectFromRawConfig(annotationsObj)

		if !reflect.DeepEqual(annotationsList, expectedAnnotations) {
			t.Fatalf("[Check 3] The processing result of the annotationsList extraction is not as expected, want %s, but got %s",
				Green(expectedAnnotations), Yellow(annotationsList))
		}
	})

	t.Logf("All processing results of the GetNestedObjectFromRawConfig method meets expectation")
}

func TestStateManagementFunc_RefreshObjectParamOriginValues(t *testing.T) {
	// 测试深度复制功能
	original := map[string]interface{}{
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
		},
		"other_param": "should_not_change",
	}

	// 测试深度复制
	copied := deepCopyInterface(original)

	// 验证深度复制是否成功
	if !reflect.DeepEqual(original, copied) {
		t.Fatalf("[Check 1] Deep copy failed, original and copied are not equal")
	}

	// 修改复制后的数据，验证原数据不受影响
	copiedMap := copied.(map[string]interface{})
	policyMap := copiedMap["policy"].(map[string]interface{})
	annotations := policyMap["annotations"].([]interface{})
	firstAnnotation := annotations[0].(map[string]interface{})
	firstAnnotation["key1"] = "modified_value"

	// 验证原数据没有被修改
	originalPolicyMap := original["policy"].(map[string]interface{})
	originalAnnotations := originalPolicyMap["annotations"].([]interface{})
	originalFirstAnnotation := originalAnnotations[0].(map[string]interface{})

	if originalFirstAnnotation["key1"] != "value1" {
		t.Fatalf("[Check 2] Original data was affected by modification of copied data")
	}

	// 测试嵌套结构安全更新
	testCases := []struct {
		name     string
		current  interface{}
		parts    []string
		value    interface{}
		expected interface{}
	}{
		{
			name: "update map nested key",
			current: map[string]interface{}{
				"policy": map[string]interface{}{
					"config": map[string]interface{}{
						"setting1": "old_value",
					},
				},
			},
			parts: []string{"policy", "config", "setting1"},
			value: "new_value",
			expected: map[string]interface{}{
				"policy": map[string]interface{}{
					"config": map[string]interface{}{
						"setting1": "new_value",
					},
				},
			},
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
		},
	}

	for i, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			currentCopy := deepCopyInterface(tt.current)
			result, err := updateNestedStructureSafely(currentCopy, tt.parts, tt.value)

			if err != nil {
				t.Fatalf("[Check %d] updateNestedStructureSafely failed: %v", i+1, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("[Check %d] Result is not as expected, want %s, but got %s",
					i+1, Green(tt.expected), Yellow(result))
			}
		})
	}

	t.Logf("All RefreshObjectParamOriginValues safety tests passed")
}
