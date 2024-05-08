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

package sdktime

import (
	"fmt"
	"strings"
	"time"
)

type SdkTime time.Time

func (t *SdkTime) UnmarshalJSON(data []byte) error {
	tmp := strings.Trim(string(data[:]), "\"")

	now, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(`2006-01-02T15:04:05`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(`2006-01-02 15:04:05`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(`2006-01-02T15:04:05+08:00`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(time.RFC3339, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(time.RFC3339Nano, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	return err
}

func (t SdkTime) MarshalJSON() ([]byte, error) {
	rs := []byte(fmt.Sprintf(`"%s"`, t.String()))
	return rs, nil
}

func (t SdkTime) String() string {
	return time.Time(t).Format(`2006-01-02T15:04:05Z`)
}
