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
	"encoding/hex"
	"net/http"
)

const (
	sdkHmacSm3     = "SDK-HMAC-SM3"
	xSdkContentSm3 = "X-Sdk-Content-Sm3"
)

type SM3Signer struct {
}

func (s SM3Signer) Sign(req *http.Request, ak, sk string) (map[string]string, error) {
	err := checkAKSK(ak, sk)
	if err != nil {
		return nil, err
	}

	processContentHeader(req, xSdkContentSm3)
	t := extractTime(req.Header.Get(HeaderXDateTime))
	req.Header.Set(HeaderXDateTime, t.UTC().Format(BasicDateFormat))
	headerDate := t.UTC().Format(BasicDateFormat)
	additionalHeaders := map[string]string{HeaderXDate: headerDate}

	signedHeaders := extractSignedHeaders(req)

	cr, err := canonicalRequest(req, signedHeaders, xSdkContentSm3, sm3HasherInst)
	if err != nil {
		return nil, err
	}

	sts, err := stringToSign(sdkHmacSm3, cr, t, sm3HasherInst)
	if err != nil {
		return nil, err
	}

	sig, err := s.signStringToSign(sts, []byte(sk))
	if err != nil {
		return nil, err
	}

	additionalHeaders[HeaderAuthorization] = authHeaderValue(sdkHmacSm3, sig, ak, signedHeaders)
	req.Header.Set(HeaderAuthorization, authHeaderValue(sdkHmacSm3, sig, ak, signedHeaders))

	return additionalHeaders, nil
}

func (SM3Signer) signStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hmac, err := sm3HasherInst.hmac([]byte(stringToSign), signingKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hmac), nil
}
