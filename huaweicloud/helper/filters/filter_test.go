package filters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter1(t *testing.T) {
	data := map[string]interface{}{
		"foo":    "bar",
		"strArr": []string{"a", "b"},
		"arr": []map[string]interface{}{
			{
				"str":     "bar1",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
			{
				"str":     "test string",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
		},
		"node": map[string]interface{}{
			"subNode": []map[string]interface{}{
				{
					"str":     "bar1",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
				{
					"str":     "test string",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
			},
		},
	}

	expected := map[string]interface{}{
		"foo":    "bar",
		"strArr": []string{"a", "b"},
		"arr": []map[string]interface{}{
			{
				"str":     "bar1",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
		},
		"node": map[string]interface{}{
			"subNode": []map[string]interface{}{
				{
					"str":     "bar1",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
				{
					"str":     "test string",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
			},
		},
	}

	rst, err := New().
		Data(data).
		From("arr").
		Where("str", "=", "bar1").
		Where("strArr", "has", []string{"a1", "b2"}). // ok
		Where("strArr", "hasContains", []string{"a1", "b2x"}).
		Where("intArr", "has", 1).
		Get()

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%v", rst), fmt.Sprintf("%v", expected))
}

func TestFilter2(t *testing.T) {
	data := map[string]interface{}{
		"foo":    "bar",
		"strArr": []string{"a", "b"},
		"arr": []map[string]interface{}{
			{
				"str":     "bar1",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
			{
				"str":     "test string",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
		},
		"node": map[string]interface{}{
			"subNode": []map[string]interface{}{
				{
					"str":     "bar1",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
				{
					"str":     "test string",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
			},
		},
	}

	expected := map[string]interface{}{
		"foo":    "bar",
		"strArr": []string{"a", "b"},
		"arr": []map[string]interface{}{
			{
				"str":     "bar1",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
			{
				"str":     "test string",
				"strArr":  []string{"a1", "b2"},
				"intArr":  []int{1, 2, 3, 11},
				"boolArr": []bool{false, false, false},
			},
		},
		"node": map[string]interface{}{
			"subNode": []map[string]interface{}{
				{
					"str":     "bar1",
					"strArr":  []string{"a1", "b2"},
					"intArr":  []int{1, 2, 3, 11},
					"boolArr": []bool{false, false, false},
				},
			},
		},
	}

	rst, err := New().
		Data(data).
		From("node.subNode").
		Where("str", "=", "bar1").
		Where("strArr", "has", []string{"a1", "b2"}). // ok
		Where("strArr", "hasContains", []string{"a1", "b2x"}).
		Where("intArr", "has", 1).
		Get()

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%v", rst), fmt.Sprintf("%v", expected))
}
