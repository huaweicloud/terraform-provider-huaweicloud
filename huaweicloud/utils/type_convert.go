package utils

import (
	"log"
	"reflect"
	"strconv"
)

// returns a pointer to the bool value
func Bool(v bool) *bool {
	return &v
}

// returns a pointer to the string value
func String(v string) *string {
	return &v
}

// returns a pointer to the string value. if v is empty, return nil
func StringIgnoreEmpty(v string) *string {
	if len(v) < 1 {
		return nil
	}
	return &v
}

// Int returns a pointer to the int value
func Int(v int) *int {
	return &v
}

// Int32 returns a pointer to the int32 value
func Int32(v int32) *int32 {
	return &v
}

// Int returns a pointer to the int value. if v is empty, return nil
func IntIgnoreEmpty(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

// Int32 returns a pointer to the int32 value. if v is empty, return nil
func Int32IgnoreEmpty(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

// Int32 returns a pointer to the int32 value
func Int64IgnoreEmpty(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

// Float64 returns a pointer to the float64 value
func Float64(v float64) *float64 {
	return &v
}

// StringToInt convert the string to int, and return the pointer of int value
func StringToInt(i *string) *int {
	if i == nil || len(*i) == 0 {
		return nil
	}

	r, err := strconv.Atoi(*i)
	if err != nil {
		log.Printf("[ERROR] convert the string %q to int failed.", *i)
	}
	return &r
}

// StringToBool convert the string to boolean, and return the pointer of boolean value
func StringToBool(v interface{}) *bool {
	if v, ok := v.(string); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Printf("[ERROR] convert the string %q to boolean failed.", v)
		}

		return &b
	}

	return nil
}

// StringValue returns the string value
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// ValueIgnoreEmpty returns to the string value. if v is empty, return nil
func ValueIgnoreEmpty(v interface{}) interface{} {
	vl := reflect.ValueOf(v)
	if (vl.Kind() != reflect.Bool) && vl.IsZero() {
		return nil
	}

	if (vl.Kind() == reflect.Array || vl.Kind() == reflect.Slice) && vl.Len() == 0 {
		return nil
	}

	return v
}
