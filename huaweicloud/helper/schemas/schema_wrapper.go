package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ListItemGetterFun func(i int) any
type SetItemGetterFun func(item any) any

type ResourceDataWrapper struct {
	ResourceData *schema.ResourceData
}

func NewSchemaWrapper(d *schema.ResourceData) *ResourceDataWrapper {
	return &ResourceDataWrapper{ResourceData: d}
}

func (s *ResourceDataWrapper) Get(key string, keepZero ...bool) any {
	val, ok := s.ResourceData.GetOk(key)
	if len(keepZero) > 0 && keepZero[0] {
		return val
	}

	if !ok {
		return nil
	}

	return val
}

func (s *ResourceDataWrapper) GetE(key string, keepDefault ...bool) any {
	if len(keepDefault) > 0 && keepDefault[0] {
		return s.ResourceData.Get(key)
	}

	val, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}
	return val
}

func (s *ResourceDataWrapper) StrToSlice(key string, keepZero ...bool) []string {
	val, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	v, _ := val.(string)
	if len(keepZero) > 0 && keepZero[0] {
		return []string{v}
	}
	if v == "" {
		return nil
	}
	return []string{v}
}

func (s *ResourceDataWrapper) IntToSlice(key string, keepZero ...bool) []int {
	val, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	v, _ := val.(int)
	if len(keepZero) > 0 && keepZero[0] {
		return []int{v}
	}
	if v == 0 {
		return nil
	}
	return []int{v}
}

func (s *ResourceDataWrapper) FloatToSlice(key string, keepZero ...bool) []float64 {
	val, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	v, _ := val.(float64)
	if len(keepZero) > 0 && keepZero[0] {
		return []float64{v}
	}
	if v == 0 {
		return nil
	}
	return []float64{v}
}

func (s *ResourceDataWrapper) BoolToSlice(key string, keepZero ...bool) []bool {
	val, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	v, _ := val.(bool)
	if len(keepZero) > 0 && keepZero[0] {
		return []bool{v}
	}
	if !v {
		return nil
	}
	return []bool{v}
}

func (s *ResourceDataWrapper) ListToSlice(key string, keepZero bool, getter ListItemGetterFun) []any {
	rst := make([]any, 0)
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		if keepZero {
			return rst
		}
		return nil
	}

	arr, _ := v.([]any)
	if len(arr) == 0 && !keepZero {
		return nil
	}

	for i := range arr {
		item := getter(i)
		rst = append(rst, item)
	}
	return rst
}

func (s *ResourceDataWrapper) ListToObject(key string, getter SetItemGetterFun) any {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.([]any)
	if len(input) >= 1 {
		return getter(0)
	}

	return nil
}

func (s *ResourceDataWrapper) ListToStrSlice(key string) []string {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	arr, _ := v.([]any)
	if len(arr) == 0 {
		return nil
	}

	rst := make([]string, 0)
	for _, item := range arr {
		if newVal, ok := item.(string); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) ListToIntSlice(key string) []int {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	arr, _ := v.([]any)
	if len(arr) == 0 {
		return nil
	}

	rst := make([]int, 0)
	for _, item := range arr {
		if newVal, ok := item.(int); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) ListToFloatSlice(key string) []float64 {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	arr, _ := v.([]any)
	if len(arr) == 0 {
		return nil
	}

	rst := make([]float64, 0)
	for _, item := range arr {
		if newVal, ok := item.(float64); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) SetToSlice(key string, getter SetItemGetterFun) []any {
	rst := make([]any, 0)
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil {
		return nil
	}

	for _, item := range input.List() {
		val := getter(item)
		rst = append(rst, val)
	}
	return rst
}

func (s *ResourceDataWrapper) SetToObject(key string, getter SetItemGetterFun) any {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil {
		return nil
	}

	list := input.List()
	if len(list) >= 1 {
		return getter(list[0])
	}

	return nil
}

func (s *ResourceDataWrapper) SetToStrSlice(key string) any {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil || len(input.List()) == 0 {
		return nil
	}

	rst := make([]string, 0)
	for _, item := range input.List() {
		if newVal, ok := item.(string); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) SetToBoolSlice(key string) any {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil {
		return nil
	}

	rst := make([]bool, 0)
	for _, item := range input.List() {
		if newVal, ok := item.(bool); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) SetToIntSlice(key string) []int {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil {
		return nil
	}

	rst := make([]int, 0)
	for _, item := range input.List() {
		if newVal, ok := item.(int); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) SetToFloatSlice(key string) []float64 {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil {
		return nil
	}

	rst := make([]float64, 0)
	for _, item := range input.List() {
		if newVal, ok := item.(float64); ok {
			rst = append(rst, newVal)
		}
	}
	return rst
}

func (s *ResourceDataWrapper) MapToStrMap(key string) any {
	v, ok := s.ResourceData.GetOk(key)
	if !ok {
		return nil
	}
	return v
}

func (s *ResourceDataWrapper) MapToMap(key string, keepZero bool, convFun SetItemGetterFun) any {
	rst := make(map[string]any)
	val, ok := s.ResourceData.GetOk(key)
	if !ok {
		if keepZero {
			return rst
		}
		return nil
	}

	mapData, _ := val.(map[string]any)
	if len(mapData) == 0 && !keepZero {
		return nil
	}

	for k, val := range mapData {
		v := convFun(val)
		mapData[k] = v
	}
	return rst
}
