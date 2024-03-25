package schemas

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type Json map[string]any

func TestSliceToList(t *testing.T) {
	initArray := []any{
		Json{
			"name": "张三",
			"age":  12,
		},
		Json{
			"name": "李四",
			"age":  13,
		},
	}
	initData := Json{
		"root": Json{
			"array": initArray,
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))
	testCases := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name:     "test_1",
			key:      "root.array",
			expected: initArray,
		}, {
			name:     "test_2",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := SliceToList(jsonObj.Get(tc.key), func(val gjson.Result) any {
				return Json{
					"name": val.Get("name").String(),
					"age":  val.Get("age").Int(),
				}
			})
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(t, expVal, actVal)
		})
	}
}

func TestSliceToStrList(t *testing.T) {
	initArr := []string{"a", "b", "c"}
	initData := Json{
		"root": Json{
			"array1": initArr,
			"array2": []string{},
			"array3": []any{"a", "b", "c"},
			"array4": []int{},
			"array5": []int{1, 2, 3},
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name:     "test_1",
			key:      "root.array1",
			expected: initArr,
		}, {
			name:     "test_2",
			key:      "root.array2",
			expected: []string{},
		}, {
			name:     "test_3",
			key:      "root.array3",
			expected: initArr,
		}, {
			name:     "test_4",
			key:      "root.array4",
			expected: []string{},
		}, {
			name:     "test_5",
			key:      "root.array5",
			expected: []string{"1", "2", "3"},
		}, {
			name:     "test_6",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := SliceToStrList(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(t, expVal, actVal)
		})
	}
}

func TestSliceToIntList(t *testing.T) {
	initArr := []int64{1, 2, 3}
	initData := Json{
		"root": Json{
			"array1": initArr,
			"array2": []int{},
			"array3": []any{1, 2, 3},
			"array4": []string{},
			"array5": []string{"1", "2", "3"},
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name:     "test_1",
			key:      "root.array1",
			expected: initArr,
		}, {
			name:     "test_2",
			key:      "root.array2",
			expected: []int64{},
		}, {
			name:     "test_3",
			key:      "root.array3",
			expected: initArr,
		}, {
			name:     "test_4",
			key:      "root.array4",
			expected: []int64{},
		}, {
			name:     "test_5",
			key:      "root.array5",
			expected: initArr,
		}, {
			name:     "test_6",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := SliceToIntList(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(t, expVal, actVal)
		})
	}
}

func TestSliceToBoolList(t *testing.T) {
	initArr := []bool{true, false, true, true}
	initData := Json{
		"root": Json{
			"array1": initArr,
			"array2": []bool{},
			"array3": []any{true, false, true, true},
			"array4": []string{},
			"array5": []string{"true", "false", "true", "true"},
			"array6": []int{1, 0, 1, 1},
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name:     "test_1",
			key:      "root.array1",
			expected: initArr,
		}, {
			name:     "test_2",
			key:      "root.array2",
			expected: []bool{},
		}, {
			name:     "test_3",
			key:      "root.array3",
			expected: initArr,
		}, {
			name:     "test_4",
			key:      "root.array4",
			expected: []bool{},
		}, {
			name:     "test_5",
			key:      "root.array5",
			expected: initArr,
		}, {
			name:     "test_6",
			key:      "root.array6",
			expected: initArr,
		}, {
			name:     "test_6",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := SliceToBoolList(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(t, expVal, actVal)
		})
	}
}

func TestSliceToFloatList(t *testing.T) {
	initArr := []float64{1.1, 2.2, 3.3, 4.4}
	initArr32 := []float32{1.1, 2.2, 3.3, 4.4}
	initArr64 := []float64{3.141592653589793, 2.718281828459045, 1234.56789}
	initData := Json{
		"root": Json{
			"array1": initArr32,
			"array2": initArr64,
			"array3": []float32{},
			"array4": []float64{},
			"array5": []any{1.1, 2.2, 3.3, 4.4},
			"array6": []any{3.141592653589793, 2.718281828459045, 1234.56789},

			"array7":  []string{"1.1", "2.2", "3.3", "4.4"},
			"array8":  []string{"3.141592653589793", "2.718281828459045", "1234.56789"},
			"array9":  []int{1, 2, 3, 4},
			"array10": []int64{1, 2, 3, 4},
			"array11": []int32{1, 2, 3, 4},
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name:     "test_1",
			key:      "root.array1",
			expected: initArr,
		}, {
			name:     "test_2",
			key:      "root.array2",
			expected: initArr64,
		}, {
			name:     "test_3",
			key:      "root.array3",
			expected: []float64{},
		}, {
			name:     "test_4",
			key:      "root.array4",
			expected: []float64{},
		}, {
			name:     "test_5",
			key:      "root.array5",
			expected: initArr,
		}, {
			name:     "test_6",
			key:      "root.array6",
			expected: initArr64,
		}, {
			name:     "test_7",
			key:      "root.array7",
			expected: initArr,
		}, {
			name:     "test_8",
			key:      "root.array8",
			expected: initArr64,
		}, {
			name:     "test_9",
			key:      "root.array9",
			expected: []float64{1, 2, 3, 4},
		}, {
			name:     "test_10",
			key:      "root.array10",
			expected: []float64{1, 2, 3, 4},
		}, {
			name:     "test_11",
			key:      "root.array11",
			expected: []float64{1, 2, 3, 4},
		}, {
			name:     "test_12",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act := SliceToFloatList(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(t, expVal, actVal)
		})
	}
}

func TestObjectToList(t *testing.T) {
	initObj := Json{
		"name": "张三",
		"age":  12,
	}

	initData := Json{
		"root": Json{
			"string": "张三",
			"int":    1,
			"bool":   true,
			"float":  3.141592653589793,

			"object": initObj,
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		convFunc func(val gjson.Result) any
		expected any
	}{
		{
			name: "test_1",
			key:  "root.string",
			convFunc: func(val gjson.Result) any {
				return val.String()
			},
			expected: []any{"张三"},
		}, {
			name: "test_2",
			key:  "root.int",
			convFunc: func(val gjson.Result) any {
				return val.Int()
			},
			expected: []any{1},
		}, {
			name: "test_3",
			key:  "root.bool",
			convFunc: func(val gjson.Result) any {
				return val.Bool()
			},
			expected: []any{true},
		}, {
			name: "test_4",
			key:  "root.float",
			convFunc: func(val gjson.Result) any {
				return val.Float()
			},
			expected: []any{3.141592653589793},
		}, {
			name: "test_5",
			key:  "root.object",
			convFunc: func(val gjson.Result) any {
				return Json{
					"name": val.Get("name").String(),
					"age":  val.Get("age").Int(),
				}
			},
			expected: []any{initObj},
		}, {
			name:     "test_6",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			act := ObjectToList(jsonObj.Get(tc.key), tc.convFunc)
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(tt, expVal, actVal)
		})
	}
}

func TestMapConverter(t *testing.T) {
	strMap := map[string]any{
		"a": "abc",
		"b": "true",
		"c": "1",
		"d": 1,
		"e": 3.141592653589793,
	}

	intMap := map[string]any{
		"a": 1,
		"b": "1",
		"c": 1415926535889793,
		"d": "1415926535889793",
		"e": "a",
	}

	floatMap := map[string]any{
		"a": 3.141592653589793,
		"b": "3.141592653589793",
		"c": 1415926535889793,
		"d": "1415926535889793",
		"e": "a",
	}

	boolMap := map[string]any{
		"a": true,
		"b": "true",
		"c": false,
		"d": "false",
		"e": 1,
		"f": 0,
		"g": "1",
		"h": "",
	}

	initData := Json{
		"root": Json{
			"map_str":   strMap,
			"map_int":   intMap,
			"map_float": floatMap,
			"map_bool":  boolMap,
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		convFunc func(val gjson.Result) any
		expected any
	}{
		{
			name: "test_1",
			key:  "root.map_str",
			convFunc: func(val gjson.Result) any {
				return val.String()
			},
			expected: map[string]interface{}{
				"a": "abc",
				"b": "true",
				"c": "1",
				"d": "1",
				"e": "3.141592653589793",
			},
		}, {
			name: "test_2",
			key:  "root.map_int",
			convFunc: func(val gjson.Result) any {
				return val.Int()
			},
			expected: map[string]interface{}{
				"a": 1,
				"b": 1,
				"c": 1415926535889793,
				"d": 1415926535889793,
				"e": 0,
			},
		}, {
			name: "test_3",
			key:  "root.map_float",
			convFunc: func(val gjson.Result) any {
				return val.Float()
			},
			expected: map[string]interface{}{
				"a": 3.141592653589793,
				"b": 3.141592653589793,
				"c": 1.415926535889793e+15,
				"d": 1.415926535889793e+15,
				"e": 0,
			},
		}, {
			name: "test_4",
			key:  "root.map_bool",
			convFunc: func(val gjson.Result) any {
				return val.Bool()
			},
			expected: map[string]interface{}{
				"a": true,
				"b": true,
				"c": false,
				"d": false,
				"e": true,
				"f": false,
				"g": true,
				"h": false,
			},
		}, {
			name:     "test_5",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			act := MapConverter(jsonObj.Get(tc.key), tc.convFunc)
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(tt, expVal, actVal)
		})
	}
}

func TestMapToStrMapSchema(t *testing.T) {
	strMap := map[string]any{
		"a": "abc",
		"b": "true",
		"c": "1",
		"d": 1,
		"e": 3.141592653589793,
	}

	initData := Json{
		"root": Json{
			"map_str": strMap,
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		expected any
	}{
		{
			name: "test_1",
			key:  "root.map_str",
			expected: map[string]string{
				"a": "abc",
				"b": "true",
				"c": "1",
				"d": "1",
				"e": "3.141592653589793",
			},
		}, {
			name:     "test_2",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			act := MapToStrMap(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(tt, expVal, actVal)
		})
	}
}

func TestMapToFloatMap(t *testing.T) {
	floatMap := map[string]any{
		"a": 3.141592653589793,
		"b": "3.141592653589793",
		"c": 1415926535889793,
		"d": "1415926535889793",
		"e": "a",
	}

	initData := Json{
		"root": Json{
			"map_float": floatMap,
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		convFunc func(val gjson.Result) any
		expected any
	}{
		{
			name: "test_1",
			key:  "root.map_float",
			convFunc: func(val gjson.Result) any {
				return val.Float()
			},
			expected: map[string]float64{
				"a": 3.141592653589793,
				"b": 3.141592653589793,
				"c": 1.415926535889793e+15,
				"d": 1.415926535889793e+15,
				"e": 0,
			},
		}, {
			name:     "test_2",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			act := MapToFloatMap(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(tt, expVal, actVal)
		})
	}
}

func TestMapToIntMap(t *testing.T) {
	intMap := map[string]any{
		"a": 1,
		"b": "1",
		"c": 1415926535889793,
		"d": "1415926535889793",
		"e": "a",
	}

	initData := Json{
		"root": Json{
			"map_int": intMap,
		},
	}
	b, err := json.Marshal(initData)
	assert.NoError(t, err)
	jsonObj := gjson.Parse(string(b))

	testCases := []struct {
		name     string
		key      string
		convFunc func(val gjson.Result) any
		expected any
	}{
		{
			name: "test_1",
			key:  "root.map_int",
			convFunc: func(val gjson.Result) any {
				return val.Int()
			},
			expected: map[string]int64{
				"a": 1,
				"b": 1,
				"c": 1415926535889793,
				"d": 1415926535889793,
				"e": 0,
			},
		}, {
			name:     "test_2",
			key:      "not_exist",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			act := MapToIntMap(jsonObj.Get(tc.key))
			expVal := fmt.Sprintf("%#v", tc.expected)
			actVal := fmt.Sprintf("%#v", act)
			fmt.Println("expected:", expVal)
			fmt.Println("  actual:", actVal)
			assert.Equal(tt, expVal, actVal)
		})
	}
}
