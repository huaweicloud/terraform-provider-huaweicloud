// Copyright 2022 Huawei Technologies Co.,Ltd.
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

package retry

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
)

// Exponential 指数退避 delay = min(最大等待时间, 基础等待时间 * (2^已重试次数))
type Exponential struct {
}

func (e *Exponential) ComputeDelayBeforeNextRetry() int32 {
	return utils.Min32(MaxDelay, BaseDelay*(utils.Pow32(3, 2)))
}
