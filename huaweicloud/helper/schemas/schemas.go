package schemas

import (
	"github.com/tidwall/gjson"
)

func SliceToList(arr gjson.Result, convFun func(val gjson.Result) any) any {
	if !arr.Exists() {
		return nil
	}

	rst := make([]any, 0)
	for _, v := range arr.Array() {
		rst = append(rst, convFun(v))
	}
	return rst
}

func SliceToStrList(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make([]string, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.String())
	}
	return rst
}

func SliceToIntList(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make([]int64, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.Int())
	}
	return rst
}

func SliceToBoolList(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make([]bool, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.Bool())
	}
	return rst
}

func SliceToFloatList(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make([]float64, 0)
	for _, v := range input.Array() {
		rst = append(rst, v.Float())
	}
	return rst
}

func ObjectToList(input gjson.Result, convFun func(val gjson.Result) any) any {
	if !input.Exists() {
		return nil
	}

	return []any{convFun(input)}
}

func MapToTypeMap(input gjson.Result, convFun func(val gjson.Result) any) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]any)
	for k, v := range input.Map() {
		rst[k] = convFun(v)
	}
	return rst
}

func MapToStrTypeMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]string)
	for k, v := range input.Map() {
		rst[k] = v.String()
	}
	return rst
}

func MapToFloatTypeMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]float64)
	for k, v := range input.Map() {
		rst[k] = v.Float()
	}
	return rst
}

func MapToIntTypeMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]int64)
	for k, v := range input.Map() {
		rst[k] = v.Int()
	}
	return rst
}
