package schemas

import (
	"time"

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

func PrimSliceConverter(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	arr := make([]any, 0)
	for _, v := range input.Array() {
		arr = append(arr, v.Value())
	}

	return convArrayType(arr)
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

func MapConverter(input gjson.Result, convFun func(val gjson.Result) any) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]any)
	for k, v := range input.Map() {
		rst[k] = convFun(v)
	}
	return rst
}

func MapToStrMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]string)
	for k, v := range input.Map() {
		rst[k] = v.String()
	}
	return rst
}

func MapToFloatMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]float64)
	for k, v := range input.Map() {
		rst[k] = v.Float()
	}
	return rst
}

func MapToIntMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]int64)
	for k, v := range input.Map() {
		rst[k] = v.Int()
	}
	return rst
}

func MapToBoolMap(input gjson.Result) any {
	if !input.Exists() {
		return nil
	}

	rst := make(map[string]bool)
	for k, v := range input.Map() {
		rst[k] = v.Bool()
	}
	return rst
}

func DateFormat(date gjson.Result, source, target string) string {
	t, err := time.Parse(source, date.String())
	if err != nil {
		return ""
	}
	return t.Format(target)
}
