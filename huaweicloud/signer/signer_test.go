// Copyright 2023 Huawei Technologies Co.,Ltd.
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

package signer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ak       = "AccessKey"
	sk       = "SecretKey"
	host     = "example.huaweicloud.com"
	endpoint = "https://" + host
	path     = "/path"
)

type TestBody struct {
	Name string
	Id   int
}

type TestParam struct {
	Name, Method, Endpoint, Path string
	Body                         interface{}
	Queries                      map[string]interface{}
	Headers                      map[string]string
}

type TestCase struct {
	TestParam
	Expected string
}

var (
	testParam1 = TestParam{
		Name:     "test1",
		Method:   "GET",
		Endpoint: endpoint,
		Path:     path,
		Body:     nil,
		Queries:  map[string]interface{}{"limit": 1},
		Headers:  map[string]string{"X-Sdk-Date": "20060102T150405Z", "TEST_UNDERSCORE": "TEST_VALUE"},
	}
	testParam2 = TestParam{
		Name:     "test2",
		Method:   "POST",
		Endpoint: endpoint,
		Path:     path,
		Body:     &TestBody{Name: "test", Id: 1},
		Queries:  map[string]interface{}{"key": "value"},
		Headers:  map[string]string{"X-Sdk-Date": "20060102T150405Z", "TEST_UNDERSCORE": "TEST_VALUE", "Content-Type": "application/json"},
	}
)

func buildReqWithTestcase(tc TestCase) (*http.Request, error) {
	baseURL, err := url.Parse(tc.Endpoint)
	if err != nil {
		return nil, err
	}

	path := strings.TrimPrefix(tc.Path, "/")
	baseURL.Path, err = url.JoinPath(baseURL.Path, path)
	if err != nil {
		return nil, err
	}

	if tc.Queries != nil {
		query := baseURL.Query()
		for k, v := range tc.Queries {
			query.Add(k, strings.TrimSpace(toString(v)))
		}
		baseURL.RawQuery = query.Encode()
	}

	var bodyReader *bytes.Reader
	if tc.Body != nil {
		bodyBytes, err := json.Marshal(tc.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(tc.Method, baseURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	if tc.Headers != nil {
		for k, v := range tc.Headers {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func TestSigner_Sign(t *testing.T) {
	cases := []TestCase{
		{
			TestParam: testParam1,
			Expected: "SDK-HMAC-SHA256 Access=AccessKey, SignedHeaders=x-sdk-date, " +
				"Signature=5a2ce64c865e0e6046321c6f3d5a77ba8413eeaf355c3166c03d58d02ac79624",
		},
		{
			TestParam: testParam2,
			Expected: "SDK-HMAC-SHA256 Access=AccessKey, SignedHeaders=x-sdk-date, " +
				"Signature=b3e46ef778ad125ce1a1ca37544ba260863dde826d6fcd806feb4407f8e0c21b",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			req, err := buildReqWithTestcase(c)
			assert.Nil(t, err)
			result, err := signerInst.Sign(req, ak, sk)
			assert.Nil(t, err)
			assert.Equal(t, c.Expected, result["Authorization"])
		})
	}
}
