// Copyright 2020 Huawei Technologies Co.,Ltd.
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package converter

import (
	"fmt"
	"reflect"
	"strconv"
)

type Int32Converter struct{}

func (i Int32Converter) CovertStringToInterface(value string) (interface{}, error) {
	i64, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		return int32(i64), nil
	}
	return int32(0), err
}

func (i Int32Converter) CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error {
	v, err := i.CovertStringToInterface(value)
	if err != nil {
		return err
	}
	val, ok := v.(int32)
	if !ok {
		return fmt.Errorf("failed to convert string (%s) to int32", value)
	}

	if isPtr {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
	return nil
}

type Int64Converter struct{}

func (i Int64Converter) CovertStringToInterface(value string) (interface{}, error) {
	i64, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return i64, nil
	}
	return int64(0), err
}

func (i Int64Converter) CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error {
	v, err := i.CovertStringToInterface(value)
	if err != nil {
		return err
	}
	val, ok := v.(int64)
	if !ok {
		return fmt.Errorf("failed to convert string (%s) to int64", value)
	}

	if isPtr {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
	return nil
}

type Float32Converter struct{}

func (i Float32Converter) CovertStringToInterface(value string) (interface{}, error) {
	f64, err := strconv.ParseFloat(value, 32)
	if err == nil {
		return float32(f64), nil
	}
	return float32(0), err
}

func (i Float32Converter) CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error {
	v, err := i.CovertStringToInterface(value)
	if err != nil {
		return err
	}
	val, ok := v.(float32)
	if !ok {
		return fmt.Errorf("failed to convert string (%s) to float32", value)
	}

	if isPtr {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
	return nil
}

type Float64Converter struct{}

func (i Float64Converter) CovertStringToInterface(value string) (interface{}, error) {
	f64, err := strconv.ParseFloat(value, 32)
	if err == nil {
		return f64, nil
	}
	return float64(0), err
}

func (i Float64Converter) CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error {
	v, err := i.CovertStringToInterface(value)
	if err != nil {
		return err
	}
	val, ok := v.(float64)
	if !ok {
		return fmt.Errorf("failed to convert string (%s) to float64", value)
	}

	if isPtr {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
	return nil
}

type BooleanConverter struct{}

func (i BooleanConverter) CovertStringToInterface(value string) (interface{}, error) {
	boolean, err := strconv.ParseBool(value)
	if err == nil {
		return boolean, nil
	}
	return false, err
}

func (i BooleanConverter) CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error {
	v, err := i.CovertStringToInterface(value)
	if err != nil {
		return err
	}

	val, ok := v.(bool)
	if !ok {
		return fmt.Errorf("failed to convert string (%s) to bool", value)
	}

	if isPtr {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
	return nil
}

type StringConverter struct{}

func (i StringConverter) CovertStringToInterface(value string) (interface{}, error) {
	return value, nil
}

func (i StringConverter) CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error {
	val := value
	if isPtr {
		field.Set(reflect.ValueOf(&val))
	} else {
		field.Set(reflect.ValueOf(val))
	}
	return nil
}
