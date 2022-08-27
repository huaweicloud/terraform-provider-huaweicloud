package utils

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	matchRuleByNumLt = -1 // number value match by : less than
	matchRuleByNumEq = 0  // number value match by : equal to
	matchRuleByNumGt = 1  // number value match by : greater than
)

// FilterSliceWithField can filter the slice all through a map filter.
// If the field is a nested value, using dot(.) to split them, e.g. "SubBlock.SubField".
// If value in the map is zero, it will be ignored.
func FilterSliceWithField(all interface{}, filter map[string]interface{}) ([]interface{}, error) {
	return filterSliceWithFieldRaw(all, filter, true)
}

// FilterSliceWithZeroField can filter the slice all through a map filter.
func FilterSliceWithZeroField(all interface{}, filter map[string]interface{}) ([]interface{}, error) {
	return filterSliceWithFieldRaw(all, filter, false)
}

func filterSliceWithFieldRaw(all interface{}, filter map[string]interface{}, ignoreZero bool) ([]interface{}, error) {
	var result []interface{}
	var matched bool

	allValue := reflect.ValueOf(all)
	if allValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("options type is not a slice")
	}

	newFilter := filter
	if ignoreZero {
		for key, val := range filter {
			keyValue := reflect.ValueOf(val)
			if keyValue.IsZero() {
				log.Printf("[DEBUG] ignore zero field %s", key)
				delete(newFilter, key)
			}
		}
	}

	for i := 0; i < allValue.Len(); i++ {
		refValue := allValue.Index(i)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		if refValue.Kind() != reflect.Struct {
			return nil, fmt.Errorf("object in slice is not a struct")
		}

		matched = true
		for key, val := range newFilter {
			actual, err := getStructField(refValue, key)
			if err != nil {
				return nil, fmt.Errorf("get slice field %s failed: %s", key, err)
			}

			actualVal := reflect.ValueOf(actual)
			if actualVal.Kind() == reflect.Ptr {
				actualVal = actualVal.Elem()
			}

			if actualVal.Interface() != val {
				log.Printf("[DEBUG] can not match slice[%d] field %s: expect %v, but got %v", i, key, val, actualVal)
				matched = false
				break
			}
		}

		if matched {
			result = append(result, refValue.Interface())
		}
	}
	return result, nil
}

func getStructField(v reflect.Value, field string) (interface{}, error) {
	var subField interface{}
	var err error
	structValue := v

	parts := strings.Split(field, ".")
	for _, key := range parts {
		subField, err = getStructFieldRaw(structValue, key)
		if err != nil {
			return nil, err
		}
		structValue = reflect.ValueOf(subField)
	}
	return subField, nil
}

func getStructFieldRaw(v reflect.Value, field string) (interface{}, error) {
	if v.Kind() == reflect.Struct {
		value := reflect.Indirect(v).FieldByName(field)
		if value.IsValid() {
			return value.Interface(), nil
		}

		return nil, fmt.Errorf("reflect: can not find the field %s", field)
	}
	return nil, fmt.Errorf("reflect: Value is not a struct")
}
