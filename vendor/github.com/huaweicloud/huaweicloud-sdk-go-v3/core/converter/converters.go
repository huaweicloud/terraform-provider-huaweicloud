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
	"reflect"
	"strconv"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
)

type Converter interface {
	CovertStringToInterface(value string) (interface{}, error)
	CovertStringToPrimitiveTypeAndSetField(field reflect.Value, value string, isPtr bool) error
}

func StringConverterFactory(vType string) Converter {
	switch vType {
	case "string":
		return StringConverter{}
	case "int32":
		return Int32Converter{}
	case "int64":
		return Int64Converter{}
	case "float32":
		return Float32Converter{}
	case "float64":
		return Float32Converter{}
	case "bool":
		return BooleanConverter{}
	default:
		return nil
	}
}

func ConvertInterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value.(type) {
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64)
	case int:
		return strconv.Itoa(value.(int))
	case uint:
		return strconv.Itoa(int(value.(uint)))
	case int8:
		return strconv.Itoa(int(value.(int8)))
	case uint8:
		return strconv.Itoa(int(value.(uint8)))
	case int16:
		return strconv.Itoa(int(value.(int16)))
	case uint16:
		return strconv.Itoa(int(value.(uint16)))
	case int32:
		return strconv.Itoa(int(value.(int32)))
	case uint32:
		return strconv.Itoa(int(value.(uint32)))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case bool:
		return strconv.FormatBool(value.(bool))
	case string:
		return value.(string)
	case []byte:
		return string(value.([]byte))
	default:
		b, err := utils.Marshal(value)
		if err != nil {
			return ""
		}
		return strings.Trim(string(b[:]), "\"")
	}
}
