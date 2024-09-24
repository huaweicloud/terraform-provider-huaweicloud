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

package sdkerr

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
)

const (
	xRequestId                  = "X-Request-Id"
	code                        = "code"
	message                     = "message"
	errorCode                   = "error_code"
	errorMsg                    = "error_msg"
	encodedAuthorizationMessage = "encoded_authorization_message"
)

type CredentialsTypeError struct {
	ErrorMessage string
}

func NewCredentialsTypeError(msg string) *CredentialsTypeError {
	c := &CredentialsTypeError{
		ErrorMessage: msg,
	}
	return c
}

func (c *CredentialsTypeError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", c.ErrorMessage)
}

type ConnectionError struct {
	ErrorMessage string
}

func NewConnectionError(msg string) *ConnectionError {
	c := &ConnectionError{
		ErrorMessage: msg,
	}
	return c
}

func (c *ConnectionError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", c.ErrorMessage)
}

type RequestTimeoutError struct {
	ErrorMessage string
}

func NewRequestTimeoutError(msg string) *RequestTimeoutError {
	rt := &RequestTimeoutError{
		ErrorMessage: msg,
	}
	return rt
}

func (rt *RequestTimeoutError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", rt.ErrorMessage)
}

type errMap map[string]interface{}

func (m errMap) getStringValue(key string) string {
	var result string

	value, isExist := m[key]
	if !isExist {
		return result
	}

	if strVal, ok := value.(string); ok {
		result = strVal
	}

	return result
}

type ServiceResponseError struct {
	StatusCode                  int    `json:"status_code"`
	RequestId                   string `json:"request_id"`
	ErrorCode                   string `json:"error_code" bson:"errorCode"`
	ErrorMessage                string `json:"error_message" bson:"errorMsg"`
	EncodedAuthorizationMessage string `json:"encoded_authorization_message"`
}

func NewServiceResponseError(resp *http.Response) *ServiceResponseError {
	sr := &ServiceResponseError{}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sr.ErrorMessage = err.Error()
		return sr
	}

	dataBuf := errMap{}
	if resp.Header.Get("Content-Type") == "application/bson" {
		err = bson.Unmarshal(data, &sr)
	} else {
		err = utils.Unmarshal(data, &dataBuf)
	}
	if err != nil {
		sr.ErrorMessage = string(data)
	} else {
		processServiceResponseError(dataBuf, sr)
		if sr.ErrorMessage == "" {
			sr.ErrorMessage = string(data)
		}
	}
	sr.StatusCode = resp.StatusCode
	sr.RequestId = resp.Header.Get(xRequestId)

	return sr
}

func processServiceResponseError(m errMap, sr *ServiceResponseError) {
	if value := m.getStringValue(encodedAuthorizationMessage); value != "" {
		sr.EncodedAuthorizationMessage = value
	}

	_code := m.getStringValue(errorCode)
	msg := m.getStringValue(errorMsg)
	if _code != "" && msg != "" {
		sr.ErrorCode = _code
		sr.ErrorMessage = msg
		return
	}

	_code = m.getStringValue(code)
	msg = m.getStringValue(message)
	if _code != "" && msg != "" {
		sr.ErrorCode = _code
		sr.ErrorMessage = msg
		return
	}

	for _, v := range m {
		if val, ok := v.(map[string]interface{}); ok {
			processServiceResponseError(val, sr)
		}
	}
}

func (sr ServiceResponseError) Error() string {
	data, err := utils.Marshal(sr)
	if err != nil {
		return fmt.Sprintf("{\"ErrorMessage\": \"%s\",\"ErrorCode\": \"%s\",\"EncodedAuthorizationMessage\": \"%s\"}",
			sr.ErrorMessage, sr.ErrorCode, sr.EncodedAuthorizationMessage)
	}
	return fmt.Sprintf(string(data))
}
