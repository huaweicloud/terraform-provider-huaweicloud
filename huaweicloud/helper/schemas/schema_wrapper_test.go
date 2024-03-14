package schemas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func CreateTestStruct() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"age": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"money": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"gender": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"suke": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plane": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tank": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"beta": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpcs": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}

	return resource
}

func InitTestStruct(t *testing.T, d *schema.ResourceData) {
	sukeSet := []interface{}{
		map[string]interface{}{
			"plane": "j20",
			"tank":  "t20",
		},
		map[string]interface{}{
			"plane": "j100",
			"tank":  "t100",
		},
	}

	betaMap := map[string]interface{}{
		"beta1": "aaa",
		"beta2": "ccc",
		"beta3": "bbb",
	}

	vpcsList := []interface{}{
		map[string]interface{}{
			"id":   "aaabbbccc",
			"cidr": "192.168.0.0/24",
		},
		map[string]interface{}{
			"id":   "aaacccbbb",
			"cidr": "192.168.0.0/20",
		},
	}

	assert.NoError(t, d.Set("name", "jack"))
	assert.NoError(t, d.Set("age", 20))
	assert.NoError(t, d.Set("money", 123.45))
	assert.NoError(t, d.Set("gender", true))
	assert.NoError(t, d.Set("suke", sukeSet))
	assert.NoError(t, d.Set("beta", betaMap))
	assert.NoError(t, d.Set("vpcs", vpcsList))
}

func TestSchemaWrapperGet(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.Get("name"))
	assert.Equal(t, "", w.Get("name", true))
	assert.Equal(t, nil, w.Get("name", false))

	assert.Equal(t, nil, w.Get("age"))
	assert.Equal(t, 0, w.Get("age", true))
	assert.Equal(t, nil, w.Get("age", false))

	assert.Equal(t, nil, w.Get("money"))
	assert.Equal(t, 0.0, w.Get("money", true))
	assert.Equal(t, nil, w.Get("money", false))

	assert.Equal(t, nil, w.Get("gender"))
	assert.Equal(t, false, w.Get("gender", true))
	assert.Equal(t, nil, w.Get("gender", false))

	assert.Equal(t, nil, w.Get("suke"))
	assert.Equal(t, nil, w.Get("suke", false))
	assert.Equal(t, fmt.Sprintf("%#v", &schema.Set{}), fmt.Sprintf("%#v", w.Get("suke", true)))

	assert.Equal(t, nil, w.Get("beta"))
	assert.Equal(t, map[string]interface{}{}, w.Get("beta", true))
	assert.Equal(t, nil, w.Get("beta", false))

	assert.Equal(t, nil, w.Get("vpcs"))
	assert.Equal(t, []interface{}{}, w.Get("vpcs", true))
	assert.Equal(t, nil, w.Get("vpcs", false))

	assert.Equal(t, nil, w.Get("vpcs.0"))
	assert.Equal(t, map[string]interface{}{}, w.Get("vpcs.0", true))
	assert.Equal(t, nil, w.Get("vpcs.0", false))

	// init data
	InitTestStruct(t, d)
	assert.Equal(t, "jack", w.Get("name"))
	assert.Equal(t, "jack", w.Get("name", true))
	assert.Equal(t, "jack", w.Get("name", false))

	assert.Equal(t, 20, w.Get("age"))
	assert.Equal(t, 20, w.Get("age", true))
	assert.Equal(t, 20, w.Get("age", false))

	assert.Equal(t, 123.45, w.Get("money"))
	assert.Equal(t, 123.45, w.Get("money", true))
	assert.Equal(t, 123.45, w.Get("money", false))

	assert.Equal(t, true, w.Get("gender"))
	assert.Equal(t, true, w.Get("gender", true))
	assert.Equal(t, true, w.Get("gender", false))

	assert.NotNil(t, w.Get("suke"))
	assert.NotNil(t, w.Get("beta"))
	assert.NotNil(t, w.Get("vpcs"))
}

func TestPrimToArray(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.PrimToArray("name"))
	assert.Equal(t, []string{}, w.PrimToArray("name", true))
	assert.Equal(t, nil, w.PrimToArray("name", false))

	assert.Equal(t, nil, w.PrimToArray("age"))
	assert.Equal(t, []int{}, w.PrimToArray("age", true))
	assert.Equal(t, nil, w.PrimToArray("age", false))

	assert.Equal(t, nil, w.PrimToArray("money"))
	assert.Equal(t, []float64{}, w.PrimToArray("money", true))
	assert.Equal(t, nil, w.PrimToArray("money", false))

	assert.Equal(t, nil, w.PrimToArray("gender"))
	assert.Equal(t, []bool{}, w.PrimToArray("gender", true))
	assert.Equal(t, nil, w.PrimToArray("gender", false))

	assert.Equal(t, nil, w.PrimToArray("foo"))
	assert.Equal(t, []any{}, w.PrimToArray("foo", true))
	assert.Equal(t, nil, w.PrimToArray("foo", false))

	assert.Equal(t, nil, w.PrimToArray("vpcs"))
	assert.Equal(t, []any{}, w.PrimToArray("vpcs", true))
	assert.Equal(t, nil, w.PrimToArray("vpcs", false))

	// init data
	InitTestStruct(t, d)

	assert.Equal(t, []string{"jack"}, w.PrimToArray("name"))
	assert.Equal(t, []string{"jack"}, w.PrimToArray("name", true))
	assert.Equal(t, []string{"jack"}, w.PrimToArray("name", false))

	assert.Equal(t, []int{20}, w.PrimToArray("age"))
	assert.Equal(t, []int{20}, w.PrimToArray("age", true))
	assert.Equal(t, []int{20}, w.PrimToArray("age", false))

	assert.Equal(t, []float64{123.45}, w.PrimToArray("money"))
	assert.Equal(t, []float64{123.45}, w.PrimToArray("money", true))
	assert.Equal(t, []float64{123.45}, w.PrimToArray("money", false))

	assert.Equal(t, []bool{true}, w.PrimToArray("gender"))
	assert.Equal(t, []bool{true}, w.PrimToArray("gender", true))
	assert.Equal(t, []bool{true}, w.PrimToArray("gender", false))

	expectedResult := []interface{}{
		map[string]interface{}{"beta1": "aaa", "beta2": "ccc", "beta3": "bbb"},
	}
	assert.Equal(t, expectedResult, w.PrimToArray("beta"))
	assert.Equal(t, expectedResult, w.PrimToArray("beta", true))
	assert.Equal(t, expectedResult, w.PrimToArray("beta", false))
}

func TestListToSlice(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.ListToSlice("name", false, func(_ int) any { return d.Get("name") }))
	assert.Equal(t, nil, w.ListToSlice("beta", false, func(_ int) any { return d.Get("beta") }))
	assert.Equal(t, nil, w.ListToSlice("vpcs", false, func(i int) any { return d.Get(fmt.Sprintf("vpcs.%d", i)) }))
	assert.Equal(t, nil, w.ListToSlice("foo", false, func(i int) any { return d.Get(fmt.Sprintf("foo.%d", i)) }))

	assert.Equal(t, []interface{}{}, w.ListToSlice("name", true, func(_ int) any { return d.Get("name") }))
	assert.Equal(t, []interface{}{}, w.ListToSlice("beta", true, func(_ int) any { return d.Get("beta") }))
	assert.Equal(t, []interface{}{}, w.ListToSlice("vpcs", true, func(i int) any {
		return d.Get(fmt.Sprintf("vpcs.%d", i))
	}))
	assert.Equal(t, []interface{}{}, w.ListToSlice("foo", true, func(i int) any {
		return d.Get(fmt.Sprintf("foo.%d", i))
	}))

	// init data
	InitTestStruct(t, d)
	assert.Equal(t, nil, w.ListToSlice("name", false, func(_ int) any { return d.Get("name") }))
	assert.Equal(t, nil, w.ListToSlice("beta", false, func(_ int) any { return d.Get("beta") }))
	assert.Equal(t, []interface{}{}, w.ListToSlice("beta", true, func(_ int) any { return d.Get("beta") }))

	assert.Equal(t, []interface{}{}, w.ListToSlice("name", true, func(_ int) any { return d.Get("name") }))
	assert.Equal(t, []interface{}{}, w.ListToSlice("suke", true, func(_ int) any { return d.Get("suke") }))
	assert.Equal(t, []interface{}{}, w.ListToSlice("beta", true, func(_ int) any { return d.Get("beta") }))

	expectedResult := []interface{}{
		map[string]interface{}{"cidr": "192.168.0.0/24", "id": "aaabbbccc"},
		map[string]interface{}{"cidr": "192.168.0.0/20", "id": "aaacccbbb"},
	}
	assert.Equal(t, expectedResult, w.ListToSlice("vpcs", true, func(i int) any {
		return d.Get(fmt.Sprintf("vpcs.%d", i))
	}))
	assert.Equal(t, expectedResult, w.ListToSlice("vpcs", false, func(i int) any {
		return d.Get(fmt.Sprintf("vpcs.%d", i))
	}))
}

func TestListToObject(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.ListToObject("name", func(_ any) any { return d.Get("name") }))
	assert.Equal(t, nil, w.ListToObject("age", func(_ any) any { return d.Get("age") }))
	assert.Equal(t, nil, w.ListToObject("suke", func(item any) any { return d.Get(fmt.Sprintf("suke.%d", item)) }))
	assert.Equal(t, nil, w.ListToObject("beta", func(_ any) any { return d.Get("beta") }))
	assert.Equal(t, nil, w.ListToObject("vpcs", func(item any) any { return d.Get(fmt.Sprintf("vpcs.%d", item)) }))

	// init data
	InitTestStruct(t, d)
	expectedResult := map[string]interface{}{"cidr": "192.168.0.0/24", "id": "aaabbbccc"}

	assert.Equal(t, nil, w.ListToObject("name", func(_ any) any { return d.Get("name") }))
	assert.Equal(t, nil, w.ListToObject("age", func(_ any) any { return d.Get("age") }))
	assert.Equal(t, nil, w.ListToObject("suke", func(item any) any { return d.Get(fmt.Sprintf("suke.%d", item)) }))
	assert.Equal(t, nil, w.ListToObject("beta", func(_ any) any { return d.Get("beta") }))
	assert.Equal(t, expectedResult, w.ListToObject("vpcs", func(item any) any {
		return d.Get(fmt.Sprintf("vpcs.%d", item))
	}))
}

func TestListToArray(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, []interface{}{}, w.ListToArray("name", true))
	assert.Equal(t, nil, w.ListToArray("name", false))
	assert.Equal(t, []interface{}{}, w.ListToArray("suke", true))
	assert.Equal(t, []interface{}{}, w.ListToArray("beta", true))
	assert.Equal(t, []interface{}{}, w.ListToArray("vpcs", true))
	assert.Equal(t, nil, w.ListToArray("vpcs", false))
	assert.Equal(t, []interface{}{}, w.ListToArray("foo", true))

	// init data
	InitTestStruct(t, d)
	assert.Equal(t, []interface{}(nil), w.ListToArray("name", true))
	assert.Equal(t, nil, w.ListToArray("name", false))
	assert.Equal(t, []interface{}(nil), w.ListToArray("suke", true))
	assert.Equal(t, nil, w.ListToArray("suke", false))
	assert.Equal(t, []interface{}(nil), w.ListToArray("beta", true))

	expectedResult := []interface{}{
		map[string]interface{}{"cidr": "192.168.0.0/24", "id": "aaabbbccc"},
		map[string]interface{}{"cidr": "192.168.0.0/20", "id": "aaacccbbb"},
	}
	assert.Equal(t, expectedResult, w.ListToArray("vpcs", true))
	assert.Equal(t, expectedResult, w.ListToArray("vpcs", false))
}

func TestConvArrayType(t *testing.T) {
	array0 := []any{}
	array1 := []any{"aa", "bb"}
	array2 := []any{1, 2, 3}
	array3 := []any{1.2, 3.4}
	array4 := []any{true, false}
	array5 := []any{map[string]string{"foo": "bar"}}

	// test
	assert.Equal(t, array0, convArrayType(array0))
	assert.Equal(t, []string{"aa", "bb"}, convArrayType(array1))
	assert.Equal(t, []int{1, 2, 3}, convArrayType(array2))
	assert.Equal(t, []float64{1.2, 3.4}, convArrayType(array3))
	assert.Equal(t, []bool{true, false}, convArrayType(array4))
	assert.Equal(t, array5, convArrayType(array5))
}

func TestSetToArray(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, []interface{}{}, w.SetToArray("name", true))
	assert.Equal(t, nil, w.SetToArray("name", false))
	assert.Equal(t, []interface{}{}, w.SetToArray("suke", true))
	assert.Equal(t, []interface{}{}, w.SetToArray("beta", true))
	assert.Equal(t, []interface{}{}, w.SetToArray("vpcs", true))
	assert.Equal(t, []interface{}{}, w.SetToArray("foo", true))

	// init data
	InitTestStruct(t, d)
	assert.Equal(t, []interface{}{}, w.SetToArray("name", true))
	assert.Equal(t, nil, w.SetToArray("name", false))
	assert.Equal(t, []interface{}{}, w.SetToArray("beta", true))
	assert.Equal(t, nil, w.SetToArray("beta", false))
	assert.Equal(t, []interface{}{}, w.SetToArray("vpcs", true))

	expectedResult := []interface{}{
		map[string]interface{}{"plane": "j100", "tank": "t100"},
		map[string]interface{}{"plane": "j20", "tank": "t20"},
	}
	assert.Equal(t, expectedResult, w.SetToArray("suke", true))
	assert.Equal(t, expectedResult, w.SetToArray("suke", false))
}

func TestSetToSlice(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, []interface{}(nil), w.SetToSlice("name", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("age", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("beta", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("vpcs", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("suke", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("foo", func(item any) any { return item }))

	// init data
	InitTestStruct(t, d)

	assert.Equal(t, []interface{}(nil), w.SetToSlice("name", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("age", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("beta", func(item any) any { return item }))
	assert.Equal(t, []interface{}(nil), w.SetToSlice("vpcs", func(item any) any { return item }))

	expectedResult := []interface{}{
		map[string]interface{}{"plane": "j100", "tank": "t100"},
		map[string]interface{}{"plane": "j20", "tank": "t20"},
	}
	assert.Equal(t, expectedResult, w.SetToSlice("suke", func(item any) any { return item }))
}

func TestSetToObject(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.SetToObject("name", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("age", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("beta", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("vpcs", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("suke", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("foo", func(item any) any { return item }))

	// init data
	InitTestStruct(t, d)

	assert.Equal(t, nil, w.SetToObject("name", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("age", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("beta", func(item any) any { return item }))
	assert.Equal(t, nil, w.SetToObject("vpcs", func(item any) any { return item }))

	expectedResult := map[string]interface{}{"plane": "j100", "tank": "t100"}
	assert.Equal(t, expectedResult, w.SetToObject("suke", func(item any) any { return item }))
}

func TestMapToStrMap(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.MapToStrMap("name"))
	assert.Equal(t, nil, w.MapToStrMap("suke"))
	assert.Equal(t, nil, w.MapToStrMap("vpcs"))
	assert.Equal(t, nil, w.MapToStrMap("beta"))

	// init data
	InitTestStruct(t, d)
	assert.Equal(t, "jack", w.MapToStrMap("name"))

	expectedResult := map[string]interface{}{"beta1": "aaa", "beta2": "ccc", "beta3": "bbb"}
	assert.Equal(t, expectedResult, w.MapToStrMap("beta"))
}

func TestMapToMap(t *testing.T) {
	// create empty struct
	d := CreateTestStruct().TestResourceData()
	w := NewSchemaWrapper(d)

	// test nil
	assert.Equal(t, nil, w.MapToMap("name", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, nil, w.MapToMap("suke", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, nil, w.MapToMap("vpcs", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, nil, w.MapToMap("beta", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))

	assert.Equal(t, map[string]interface{}{}, w.MapToMap("name", true, func(item any) any {
		return fmt.Sprintf("xxx-%v", item)
	}))
	assert.Equal(t, map[string]interface{}{}, w.MapToMap("beta", true, func(item any) any {
		return fmt.Sprintf("xxx-%v", item)
	}))

	// init data
	InitTestStruct(t, d)
	assert.Equal(t, nil, w.MapToMap("name", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, nil, w.MapToMap("suke", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, nil, w.MapToMap("vpcs", false, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, map[string]interface{}{}, w.MapToMap("name", true, func(item any) any {
		return fmt.Sprintf("xxx-%v", item)
	}))
	assert.Equal(t, map[string]interface{}{}, w.MapToMap("suke", true, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))

	expectedResult := map[string]interface{}{"beta1": "xxx-aaa", "beta2": "xxx-ccc", "beta3": "xxx-bbb"}
	assert.Equal(t, expectedResult, w.MapToMap("beta", true, func(item any) any { return fmt.Sprintf("xxx-%v", item) }))
	assert.Equal(t, expectedResult, w.MapToMap("beta", false, func(item any) any {
		return fmt.Sprintf("xxx-%v", item)
	}))
}
