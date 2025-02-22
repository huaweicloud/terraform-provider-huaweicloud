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
	"math/big"
	"net/http"

	"github.com/tjfoc/gmsm/sm2"
)

const (
	sdkSm2Sm3 = "SDK-SM2-SM3"
)

var (
	curveSm2     = sm2.P256Sm2()
	sm2nMinusTwo = new(big.Int).Sub(new(big.Int).Set(curveSm2.Params().N), big.NewInt(2))
)

type SM2SM3Signer struct {
}

func (s SM2SM3Signer) Sign(req *http.Request, ak, sk string) (map[string]string, error) {
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

	sts, err := stringToSign(sdkSm2Sm3, cr, t, sm3HasherInst)
	if err != nil {
		return nil, err
	}

	signingKey, err := s.GetSigningKey(ak, sk)
	if err != nil {
		return nil, err
	}

	sig, err := signStringToSign(sts, signingKey)
	if err != nil {
		return nil, err
	}

	additionalHeaders[HeaderAuthorization] = authHeaderValue(sdkSm2Sm3, sig, ak, signedHeaders)
	req.Header.Set(HeaderAuthorization, authHeaderValue(sdkSm2Sm3, sig, ak, signedHeaders))

	return additionalHeaders, nil
}

func (s SM2SM3Signer) GetSigningKey(ak, sk string) (ISigningKey, error) {
	privateInt, err := derivePrivateInt(sdkSm2Sm3, ak, sk, sm2nMinusTwo, sm3HasherInst)
	if err != nil {
		return nil, err
	}

	return s.deriveSigningKey(privateInt), nil
}

func (SM2SM3Signer) deriveSigningKey(priv *big.Int) ISigningKey {
	privateKey := new(sm2.PrivateKey)
	privateKey.PublicKey.Curve = curveSm2
	privateKey.D = priv
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curveSm2.ScalarBaseMult(priv.Bytes())
	return SM2SigningKey{privateKey: privateKey}
}
