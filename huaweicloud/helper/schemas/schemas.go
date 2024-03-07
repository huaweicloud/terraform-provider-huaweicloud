package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ConvSchemaFun func(val any) any
type ConvResultFun func(val gjson.Result) any

func ToString(d *schema.ResourceData, key string) string {
	input, ok := d.GetOk(key)
	if !ok {
		return ""
	}
	s, _ := input.(string)
	return s
}

func ListToSlice(input any, convFun ConvSchemaFun) any {
	if input == nil {
		return nil
	}

	arr, ok := input.([]any)
	if !ok || len(arr) == 0 {
		return nil
	}

	rst := make([]any, 0)
	for _, val := range arr {
		newVal := convFun(val)
		rst = append(rst, newVal)
	}

	return rst
}

func ListToObject(input any, convFun ConvSchemaFun) any {
	if input == nil {
		return nil
	}

	arr := ListToSlice(input, convFun)
	if arr == nil {
		return nil
	}

	if tmpArr, ok := arr.([]any); ok && len(tmpArr) > 0 {
		return tmpArr[0]
	}

	return nil
}

func SetToSlice(input any, convFun ConvSchemaFun) any {
	set, _ := input.(*schema.Set)

	if set == nil {
		return nil
	}
	return ListToSlice(set.List(), convFun)
}

func MapToSlice(input any, convFun ConvSchemaFun) any {
	rst := make([]any, 0)

	arr, _ := input.(map[string]interface{})
	for _, val := range arr {
		newVal := convFun(val)
		rst = append(rst, newVal)
	}

	return rst
}

func MapToMap(input any, convFun ConvSchemaFun) any {
	rst := make(map[string]any, 0)

	mapData, _ := input.(map[string]any)
	for key, val := range mapData {
		v := convFun(val)
		mapData[key] = v
	}

	return rst
}

func MapToStrMap(input any) any {
	rst := make(map[string]any, 0)

	mapData, _ := input.(map[string]any)
	for key, val := range mapData {
		v, _ := val.(string)
		mapData[key] = v
	}

	return rst
}

func OmitEmpty(input map[string]any, keepZero ...bool) any {
	if len(keepZero) > 0 && keepZero[0] {
		return input
	}

	newInput := utils.RemoveNil(input)
	if len(newInput) == 0 {
		return nil
	}

	return input
}

func TypeListToStringSlice(input any) any {
	arr := ListToSlice(input, func(val any) any {
		str, _ := val.(string)
		return str
	})
	tmpArr, ok := arr.([]any)
	if !ok {
		return nil
	}

	rst := make([]string, 0)
	for _, s := range tmpArr {
		if str, ok := s.(string); ok {
			rst = append(rst, str)
		}
	}
	return rst
}

func TypeListToIntSlice(input any) any {
	arr := ListToSlice(input, func(val any) any {
		str, _ := val.(int)
		return str
	})
	rst, _ := arr.([]int)
	return rst
}

func TypeListToFloat64Slice(input any) any {
	arr := ListToSlice(input, func(val any) any {
		str, _ := val.(float64)
		return str
	})
	rst, _ := arr.([]float64)
	return rst
}

func TypeSetToStringSlice(input any) any {
	if input == nil {
		return nil
	}

	set, _ := input.(*schema.Set)
	if set == nil {
		return nil
	}
	return TypeListToStringSlice(set.List())
}

func TypeSetToIntSlice(input any) any {
	set, _ := input.(*schema.Set)
	if set == nil {
		return nil
	}
	return TypeListToIntSlice(set.List())
}

func TypeSetToFloat64Slice(input any) any {
	set, _ := input.(*schema.Set)
	if set == nil {
		return nil
	}
	return TypeListToFloat64Slice(set.List())
}

func SliceToList(arr gjson.Result, convFun ConvResultFun) []any {
	if !arr.Exists() {
		return nil
	}

	rst := make([]any, 0)
	for _, r := range arr.Array() {
		v := convFun(r)
		rst = append(rst, v)
	}
	return rst
}

func SliceToStrList(input gjson.Result) []string {
	if !input.Exists() {
		return nil
	}

	rst := make([]string, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.String())
	}
	return rst
}

func SliceToIntList(input gjson.Result) []int64 {
	if !input.Exists() {
		return nil
	}

	rst := make([]int64, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.Int())
	}
	return rst
}

func SliceToBoolList(input gjson.Result) []bool {
	if !input.Exists() {
		return nil
	}

	rst := make([]bool, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.Bool())
	}
	return rst
}

func SliceToFloatList(input gjson.Result) []float64 {
	if !input.Exists() {
		return nil
	}

	rst := make([]float64, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.Float())
	}
	return rst
}

func ObjectToList(input gjson.Result, convFun ConvResultFun) []any {
	if !input.Exists() {
		return nil
	}

	return []any{convFun(input)}
}

func MapToTypeMap(input gjson.Result, convFun ConvResultFun) map[string]any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]any, 0)
	for k, r := range input.Map() {
		v := convFun(r)
		rst[k] = v
	}
	return rst
}

func MapToStrTypeMap(input gjson.Result) map[string]string {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]string, 0)
	for k, r := range input.Map() {
		v := r.String()
		rst[k] = v
	}
	return rst
}

func MapToFloatTypeMap(input gjson.Result) map[string]float64 {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]float64, 0)
	for k, r := range input.Map() {
		v := r.Float()
		rst[k] = v
	}
	return rst
}

func MapToIntTypeMap(input gjson.Result) map[string]int64 {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]int64, 0)
	for k, r := range input.Map() {
		v := r.Int()
		rst[k] = v
	}
	return rst
}
