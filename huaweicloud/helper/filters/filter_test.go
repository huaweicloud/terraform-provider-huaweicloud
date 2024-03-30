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

func TestFilterData(t *testing.T) {
	filter := New()
	data1 := map[string]interface{}{
		"foo1": "bar1",
	}
	data2 := map[string]interface{}{
		"foo2": "bar2",
	}

	filter.Data(data1)
	assert.NotNil(t, filter.jsonData)
	assert.Equal(t, data1, filter.jsonData)

	filter.Data(data2)
	assert.Equal(t, data2, filter.jsonData)
}

func TestFilterFrom(t *testing.T) {
	filter := New()
	node1 := "node1"
	node2 := "node2"

	filter.From(node1)
	assert.NotNil(t, filter.node)
	assert.Equal(t, node1, filter.node)

	filter.From(node2)
	assert.Equal(t, node2, filter.node)
}

func TestFilterGetFrom(t *testing.T) {
	filter := New()
	node1 := "node1"
	node2 := "node2"

	filter.From(node1)
	assert.NotNil(t, filter.GetFrom())
	assert.Equal(t, node1, filter.GetFrom())

	filter.From(node2)
	assert.Equal(t, node2, filter.GetFrom())
}

func TestFilterGetQueries(t *testing.T) {
	filter := New()
	queriesData1 := []QueryCond{
		{
			Key:      "",
			Operator: "",
			Value:    nil,
		},
	}
	queriesData2 := []QueryCond{
		{
			Key:      "foo1",
			Operator: "op1",
			Value:    "bar1",
		},
		{
			Key:      "foo2",
			Operator: "op2",
			Value:    2,
		},
	}
	filter.queries = append(filter.queries, queriesData1...)
	assert.NotNil(t, filter.GetQueries())
	assert.Equal(t, queriesData1, filter.GetQueries())

	filter.queries = append(filter.queries, queriesData2...)
	queriesData1 = append(queriesData1, queriesData2...)
	assert.NotNil(t, filter.GetQueries())
	assert.Equal(t, queriesData1, filter.GetQueries())
}

func TestFilterWhere(t *testing.T) {
	filter := New()
	key1, operator1, value1 := "foo1", ">", "bar1"
	key2, operator2, value2 := "foo2", "<", "bar2"
	key3, operator3 := "foo3", "="

	filter.Where(key1, operator1, value1)
	assert.Equal(t, 1, len(filter.queries))
	assert.Equal(t, key1, filter.queries[0].Key)
	assert.Equal(t, operator1, filter.queries[0].Operator)
	assert.Equal(t, value1, filter.queries[0].Value)

	filter.Where(key2, operator2, value2)
	assert.Equal(t, 2, len(filter.queries))
	assert.Equal(t, key2, filter.queries[1].Key)
	assert.Equal(t, operator2, filter.queries[1].Operator)
	assert.Equal(t, value2, filter.queries[1].Value)

	filter.Where(key3, operator3, nil)
	assert.Equal(t, 2, len(filter.queries))
	assert.Equal(t, key2, filter.queries[1].Key)
	assert.Equal(t, operator2, filter.queries[1].Operator)
	assert.Equal(t, value2, filter.queries[1].Value)
}

func TestFilterSlice(t *testing.T) {
	filter := New()
	data := []map[string]interface{}{
		{
			"name": "Tom",
			"age":  20,
		},
		{
			"name": "Jerry",
			"age":  10,
		},
		{
			"name": "",
			"age":  nil,
		},
	}
	filter.Data(data)
	result, err := filter.Get()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, filter.jsonData)
}

func TestFilterJson(t *testing.T) {
	filter := New()
	data := map[string]interface{}{
		"user": []map[string]interface{}{
			{
				"name": "Tom",
				"age":  20,
			},
			{
				"name": "Jerry",
				"age":  10,
			},
			{
				"name": "",
				"age":  nil,
			},
		},
	}
	queriesData1 := []QueryCond{
		{
			Key:      "",
			Operator: "",
			Value:    nil,
		},
	}
	// test node is empty
	filter.Data(data)
	result, err := filter.Get()
	assert.Error(t, err)
	assert.Nil(t, result)

	// test node is not empty
	filter.From("users")
	result, err = filter.Get()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, filter.jsonData)

	// test queries is not empty
	filter.queries = append(filter.queries, queriesData1...)
	result, err = filter.Get()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestToStrSlice(t *testing.T) {
	// test integer array to string slice
	var array1 [5]string
	array2 := [5]int{1, 3, 2, 4, 5}
	expectedArr := []string{"1", "3", "2", "4", "5"}

	arrResult1 := toStrSlice(array1)
	arrResult2 := toStrSlice(array2)
	assert.Equal(t, []string{"", "", "", "", ""}, arrResult1)
	assert.Equal(t, expectedArr, arrResult2)

	// test integer slice to string slice
	var slice1 []string
	slice2 := []int{1, 3, 2, 4, 5}
	expectedSlice := []string{"1", "3", "2", "4", "5"}

	sliceResult1 := toStrSlice(slice1)
	sliceResult2 := toStrSlice(slice2)
	assert.Equal(t, []string{}, sliceResult1)
	assert.Equal(t, expectedSlice, sliceResult2)

	// test other data
	var expectedResult []string
	var map1 map[string]string
	map2 := map[string]string{"foo1": "bar1", "foo2": "bar2"}

	mapResult1 := toStrSlice(map1)
	mapResult2 := toStrSlice(map2)
	assert.Equal(t, expectedResult, mapResult1)
	assert.Equal(t, expectedResult, mapResult2)
}

func TestContains(t *testing.T) {
	strSlice := []string{"Tom", "Jerry"}
	assert.True(t, contains(strSlice, "Tom"))
	assert.False(t, contains(strSlice, "xxx"))
	assert.False(t, contains(strSlice, ""))

	intSlice := []int{1, 2, 3, 4, 5}
	assert.True(t, contains(intSlice, 1))
	assert.False(t, contains(intSlice, 10))

	var emptySlice []int
	assert.False(t, contains(emptySlice, 10))
}

func TestHas(t *testing.T) {
	x1 := []string{"Tom", "Jerry", "Suke", "Beta"}
	y1 := []string{"Tom"}          // a y in x
	y2 := []string{"Tom", "Jerry"} // many y in x
	y3 := []string{"foo"}          // y not in x
	y4 := []string{"Tom", "foo"}   // part y in x

	x2 := "Tom"
	y5 := "Tom"
	y6 := "Jerry"

	result1, err := has(x1, y1)
	assert.NoError(t, err)
	assert.Equal(t, true, result1)

	result2, err := has(x1, y2)
	assert.NoError(t, err)
	assert.Equal(t, true, result2)

	result3, err := has(x1, y3)
	assert.NoError(t, err)
	assert.Equal(t, false, result3)

	result4, err := has(x1, y4)
	assert.NoError(t, err)
	assert.Equal(t, true, result4)

	result5, err := has(x2, y5)
	assert.NoError(t, err)
	assert.Equal(t, false, result5)

	result6, err := has(x2, y6)
	assert.NoError(t, err)
	assert.Equal(t, false, result6)
}

func TestHasContain(t *testing.T) {
	x1 := []string{"Tom", "Jerry", "Suke", "Beta"}
	y1 := []string{"Tom"}        // a y in x
	y2 := []string{"Tom", "Jer"} // many y in x
	y3 := []string{"foo"}        // y not in x
	y4 := []string{"Tom", "foo"} // part y in x

	x2 := "Tom"
	y5 := "Tom"
	y6 := "Jerry"

	result1, err := hasContain(x1, y1)
	assert.NoError(t, err)
	assert.Equal(t, true, result1)

	result2, err := hasContain(x1, y2)
	assert.NoError(t, err)
	assert.Equal(t, true, result2)

	result3, err := hasContain(x1, y3)
	assert.NoError(t, err)
	assert.Equal(t, false, result3)

	result4, err := hasContain(x1, y4)
	assert.NoError(t, err)
	assert.Equal(t, true, result4)

	result5, err := hasContain(x2, y5)
	assert.NoError(t, err)
	assert.Equal(t, false, result5)

	result6, err := hasContain(x2, y6)
	assert.NoError(t, err)
	assert.Equal(t, false, result6)
}
