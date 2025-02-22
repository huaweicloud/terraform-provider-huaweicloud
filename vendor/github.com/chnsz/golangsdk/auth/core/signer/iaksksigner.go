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
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	HmacSHA256      = "HmacSHA256"
	HmacSM3         = "HmacSM3"
	EcdsaP256SHA256 = "EcdsaP256SHA256"
	SM2SM3          = "SM2SM3"

	xSdkContentSha256 = "X-Sdk-Content-Sha256"

	BasicDateFormat     = "20060102T150405Z"
	HeaderXDate         = "X-Sdk-Date"
	HeaderHost          = "host"
	HeaderAuthorization = "Authorization"
)

type SigningAlgorithm string

type IAKSKSigner interface {
	Sign(req *http.Request, ak, sk string) (map[string]string, error)
}

func GetSigner(alg SigningAlgorithm) (IAKSKSigner, error) {
	switch alg {
	case HmacSM3:
		return sm3SignerInst, nil
	case EcdsaP256SHA256:
		return p256sha256SignerInst, nil
	case SM2SM3:
		return sm2sm3SignerInst, nil
	default:
		return nil, errors.New("unsupported signing algorithm: " + string(alg))
	}
}

func checkAKSK(ak, sk string) error {
	if ak == "" {
		return errors.New("ak is required in credentials")
	}
	if sk == "" {
		return errors.New("sk is required in credentials")
	}

	return nil
}

// stringToSign Create a "String to Sign".
func stringToSign(alg, canonicalRequest string, t time.Time, hasher iHasher) (string, error) {
	canonicalRequestHash, err := hasher.hashHexString([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s", alg, t.UTC().Format(BasicDateFormat), canonicalRequestHash), nil
}

// authHeaderValue Get the finalized value for the "Authorization" header.
// The signature parameter is the output from stringToSign
func authHeaderValue(alg, sig, ak string, signedHeaders []string) string {
	return fmt.Sprintf("%s Access=%s, SignedHeaders=%s, Signature=%s",
		alg,
		ak,
		strings.Join(signedHeaders, ";"),
		sig)
}

func processContentHeader(req *http.Request, contentHeader string) {
	contentType := req.Header.Get("Content-Type")

	if contentType != "" && !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "application/bson") {
		req.Header.Set(contentHeader, "UNSIGNED-PAYLOAD")
	}
}

func canonicalRequest(req *http.Request, signedHeaders []string, contentHeader string, hasher iHasher) (string, error) {
	hexEncode, err := getContentHash(req, contentHeader, hasher)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method,
		canonicalURI(req),
		canonicalQueryString(req),
		canonicalHeaders(req, signedHeaders),
		strings.Join(signedHeaders, ";"), hexEncode), nil
}

func getContentHash(req *http.Request, contentHeader string, hasher iHasher) (string, error) {
	if content := req.Header.Get(contentHeader); content != "" {
		return content, nil
	}

	buffer, err := getBodyToBytes(req)
	if err != nil {
		return "", err
	}

	data := buffer.Bytes()
	hexEncode, err := hasher.hashHexString(data)
	if err != nil {
		return "", err
	}
	return hexEncode, nil
}

func extractTime(date string) time.Time {
	t, err := time.Parse(BasicDateFormat, date)
	if date == "" || err != nil {
		return time.Now()
	}
	return t
}

// canonicalURI returns request uri
func canonicalURI(r *http.Request) string {
	pattens := strings.Split(r.URL.Path, "/")

	var uri []string
	for _, v := range pattens {
		uri = append(uri, escape(v))
	}

	urlPath := strings.Join(uri, "/")
	if len(urlPath) == 0 || urlPath[len(urlPath)-1] != '/' {
		urlPath += "/"
	}

	return urlPath
}

func canonicalQueryString(r *http.Request) string {
	var keys []string
	query := r.URL.Query()
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var a []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(query[key])
		for _, v := range query[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			a = append(a, kv)
		}
	}
	queryStr := strings.Join(a, "&")
	r.URL.RawQuery = queryStr

	return queryStr
}

func canonicalHeaders(r *http.Request, signerHeaders []string) string {
	var a []string
	header := make(map[string][]string)
	for k, v := range r.Header {
		header[strings.ToLower(k)] = v
	}

	for _, key := range signerHeaders {
		value := header[key]
		if strings.EqualFold(key, HeaderHost) {
			header[HeaderHost] = []string{r.Host}
		}

		sort.Strings(value)
		for _, v := range value {
			a = append(a, key+":"+strings.TrimSpace(v))
		}
	}

	return fmt.Sprintf("%s\n", strings.Join(a, "\n"))
}

func extractSignedHeaders(r *http.Request) []string {
	var sh []string
	for key := range r.Header {
		if strings.HasPrefix(strings.ToLower(key), "content-type") || strings.Contains(key, "_") {
			continue
		}
		sh = append(sh, strings.ToLower(key))
	}
	sort.Strings(sh)

	return sh
}

func getBodyToBytes(req *http.Request) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		defer req.Body.Close()

		contentType := req.Header.Get("Content-Type")

		switch contentType {
		case "application/xml":
			var data interface{}
			if err := xml.Unmarshal(bodyBytes, &data); err != nil {
				return nil, err
			}
			encoder := xml.NewEncoder(buf)
			if err := encoder.Encode(data); err != nil {
				return nil, err
			}

		case "application/bson":
			var data interface{}
			if err := bson.Unmarshal(bodyBytes, &data); err != nil {
				return nil, err
			}
			bsonData, err := bson.Marshal(data)
			if err != nil {
				return nil, err
			}
			buf.Write(bsonData)

		default:
			var data interface{}
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return nil, err
			}
			encoder := json.NewEncoder(buf)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(data); err != nil {
				return nil, err
			}
		}
	}

	return buf, nil
}
