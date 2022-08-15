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

package signer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"golang.org/x/crypto/hkdf"
	"io"
	"strings"
	"time"
)

const (
	DerivationAlgorithm = "V11-HMAC-SHA256"
	DerivedDateFormat   = "20060102"
)

// DerivationAuthHeaderValue Get the finalized value for the "Authorization" header. The signature parameter is the output from SignStringToSign
func DerivationAuthHeaderValue(signature, accessKey string, info string, signedHeaders []string) string {
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", DerivationAlgorithm, accessKey, info, strings.Join(signedHeaders, ";"), signature)
}

// Get the derivation key for derived credential.
func GetDerivationKey(accessKey string, secretKey string, info string) (string, error) {
	hash := sha256.New
	derivationKeyReader := hkdf.New(hash, []byte(secretKey), []byte(accessKey), []byte(info))
	derivationKey := make([]byte, 32)
	_, err := io.ReadFull(derivationKeyReader, derivationKey)
	return hex.EncodeToString(derivationKey), err
}

// StringToSignDerived Create a "String to Sign".
func StringToSignDerived(canonicalRequest string, info string, t time.Time) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n%s\n%s\n%x",
		DerivationAlgorithm, t.UTC().Format(BasicDateFormat), info, hash.Sum(nil)), nil
}

// SignDerived SignRequest set Authorization header
func SignDerived(r *request.DefaultHttpRequest, ak string, sk string, derivedAuthServiceName string, regionId string) (map[string]string, error) {
	if derivedAuthServiceName == "" || regionId == "" {
		panic("DerivedAuthServiceName and RegionId in credential is required when using derived auth")
	}

	var err error
	var t time.Time
	var headerParams = make(map[string]string)
	userHeaders := r.GetHeaderParams()
	if date, ok := userHeaders[HeaderXDate]; ok {
		t, err = time.Parse(BasicDateFormat, date)
		if date == "" || err != nil {
			t = time.Now()
			userHeaders[HeaderXDate] = t.UTC().Format(BasicDateFormat)
			headerParams[HeaderXDate] = t.UTC().Format(BasicDateFormat)
		}
	} else {
		t = time.Now()
		userHeaders[HeaderXDate] = t.UTC().Format(BasicDateFormat)
		headerParams[HeaderXDate] = t.UTC().Format(BasicDateFormat)
	}
	signedHeaders := SignedHeaders(userHeaders)
	canonicalRequest, err := CanonicalRequest(r, signedHeaders)
	if err != nil {
		return nil, err
	}
	info := t.UTC().Format(DerivedDateFormat) + "/" + regionId + "/" + derivedAuthServiceName
	stringToSign, err := StringToSignDerived(canonicalRequest, info, t)
	if err != nil {
		return nil, err
	}
	derivedSk, err := GetDerivationKey(ak, sk, info)
	if err != nil {
		return nil, err
	}
	signature, err := SignStringToSign(stringToSign, []byte(derivedSk))
	if err != nil {
		return nil, err
	}
	headerParams[HeaderAuthorization] = DerivationAuthHeaderValue(signature, ak, info, signedHeaders)
	return headerParams, nil
}
