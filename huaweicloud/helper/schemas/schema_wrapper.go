package schemas

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type ListItemGetterFun func(i int) any
type SetItemGetterFun func(item any) any

type ResourceDataWrapper struct {
	*schema.ResourceData
}

func NewSchemaWrapper(d *schema.ResourceData) *ResourceDataWrapper {
	return &ResourceDataWrapper{ResourceData: d}
}

func (s *ResourceDataWrapper) NewClient(cfg *config.Config, serviceName string) (*golangsdk.ServiceClient, error) {
	region := cfg.GetRegion(s.ResourceData)
	client, err := cfg.NewServiceClient(serviceName, region)
	if err != nil {
		return nil, fmt.Errorf("error creating %s client: %s", serviceName, err)
	}

	return client, nil
}

func (s *ResourceDataWrapper) Get(key string, keepZero ...bool) any {
	val, ok := s.GetOk(key)
	if !ok && s.keepZero(keepZero) {
		return val
	}

	if !ok {
		return nil
	}

	return val
}

func (s *ResourceDataWrapper) GetToInt(key string, keepZero ...bool) any {
	val, ok := s.GetOk(key)
	if !ok && s.keepZero(keepZero) {
		return 0
	}

	if !ok {
		return nil
	}

	str, _ := val.(string)
	intVal, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		// lintignore:R009
		panic(fmt.Sprintf(`* "%s": "%v", value is incorrect, it should be a numeric string, for example: "12"`, key, val))
	}

	return intVal
}

func (s *ResourceDataWrapper) GetToBool(key string, keepZero ...bool) any {
	val, ok := s.GetOk(key)
	if !ok && s.keepZero(keepZero) {
		return 0
	}

	if !ok {
		return nil
	}

	str, _ := val.(string)
	boolVal, err := strconv.ParseBool(str)

	if err != nil {
		// lintignore:R009
		panic(fmt.Sprintf(`* "%s": "%v", value is incorrect, it should be bool type string, for example: "false"`, key, val))
	}

	return boolVal
}

// PrimToArray primitive data type to Array
func (s *ResourceDataWrapper) PrimToArray(key string, keepZero ...bool) any {
	val, ok := s.GetOk(key)
	switch vv := val.(type) {
	case string:
		if !ok {
			if s.keepZero(keepZero) {
				return []string{}
			}
			return nil
		}

		return []string{vv}
	case int:
		if !ok {
			if s.keepZero(keepZero) {
				return []int{}
			}
			return nil
		}
		return []int{vv}
	case float64:
		if !ok {
			if s.keepZero(keepZero) {
				return []float64{}
			}
			return nil
		}
		return []float64{vv}
	case bool:
		if !ok {
			if s.keepZero(keepZero) {
				return []bool{}
			}
			return nil
		}
		return []bool{vv}
	}

	if !ok {
		if s.keepZero(keepZero) {
			return []any{}
		}
		return nil
	}
	return []any{val}
}

func (s *ResourceDataWrapper) ListToSlice(key string, keepZero bool, getter ListItemGetterFun) any {
	rst := make([]any, 0)
	v, ok := s.GetOk(key)
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
	v, ok := s.GetOk(key)
	if !ok {
		return nil
	}

	input, _ := v.([]any)
	if len(input) >= 1 {
		return getter(0)
	}

	return nil
}

// ListToArray convert schema.TypeList to Array
func (s *ResourceDataWrapper) ListToArray(key string, keepZero ...bool) any {
	v, ok := s.GetOk(key)
	if !ok {
		if s.keepZero(keepZero) {
			return make([]any, 0)
		}
		return nil
	}
	arr, _ := v.([]any)
	if len(arr) == 0 && !s.keepZero(keepZero) {
		return nil
	}
	return convArrayType(arr)
}

// SetToArray convert schema.Set to Array
func (s *ResourceDataWrapper) SetToArray(key string, keepZero ...bool) any {
	v, ok := s.GetOk(key)
	if !ok {
		if s.keepZero(keepZero) {
			return []any{}
		}
		return nil
	}

	input, _ := v.(*schema.Set)
	if input == nil && s.keepZero(keepZero) {
		return []any{}
	}
	if input == nil {
		return nil
	}

	arr := input.List()
	if len(arr) == 0 && !s.keepZero(keepZero) {
		return nil
	}
	return convArrayType(arr)
}

func (s *ResourceDataWrapper) SetToSlice(key string, getter SetItemGetterFun) []any {
	rst := make([]any, 0)
	v, ok := s.GetOk(key)
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

func (s *ResourceDataWrapper) MapToStrMap(key string) any {
	v, ok := s.GetOk(key)
	if !ok {
		return nil
	}
	return v
}

func (s *ResourceDataWrapper) MapToMap(key string, keepZero bool, convFun SetItemGetterFun) any {
	rst := make(map[string]any)
	val, ok := s.GetOk(key)
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

	for k, v := range mapData {
		rst[k] = convFun(v)
	}
	return rst
}

func (*ResourceDataWrapper) keepZero(keepZero []bool) bool {
	if len(keepZero) == 0 {
		return false
	}
	return keepZero[0]
}

func convArrayType(arr []any) any {
	if len(arr) == 0 {
		return arr
	}

	switch arr[0].(type) {
	case string:
		rst := make([]string, 0)
		for _, v := range arr {
			vv, _ := v.(string)
			rst = append(rst, vv)
		}
		return rst
	case int:
		rst := make([]int, 0)
		for _, v := range arr {
			vv, _ := v.(int)
			rst = append(rst, vv)
		}
		return rst
	case int64:
		rst := make([]int64, 0)
		for _, v := range arr {
			vv, _ := v.(int64)
			rst = append(rst, vv)
		}
		return rst
	case float32:
		rst := make([]float32, 0)
		for _, v := range arr {
			vv, _ := v.(float32)
			rst = append(rst, vv)
		}
		return rst
	case float64:
		rst := make([]float64, 0)
		for _, v := range arr {
			vv, _ := v.(float64)
			rst = append(rst, vv)
		}
		return rst
	case bool:
		rst := make([]bool, 0)
		for _, v := range arr {
			vv, _ := v.(bool)
			rst = append(rst, vv)
		}
		return rst
	}

	return arr
}
