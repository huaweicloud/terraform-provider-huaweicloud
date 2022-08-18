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

package impl

import (
	"bytes"
	"crypto/tls"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/exchange"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/response"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type DefaultHttpClient struct {
	httpHandler  *httphandler.HttpHandler
	httpConfig   *config.HttpConfig
	transport    *http.Transport
	goHttpClient *http.Client
}

func NewDefaultHttpClient(httpConfig *config.HttpConfig) *DefaultHttpClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: httpConfig.IgnoreSSLVerification},
	}

	if httpConfig.DialContext != nil {
		transport.DialContext = httpConfig.DialContext
	}

	if httpConfig.HttpProxy != nil {
		proxyUrl := httpConfig.HttpProxy.GetProxyUrl()
		if proxyUrl != "" {
			proxy, _ := url.Parse(proxyUrl)
			transport.Proxy = http.ProxyURL(proxy)
		}
	}

	client := &DefaultHttpClient{
		transport:  transport,
		httpConfig: httpConfig,
	}

	client.goHttpClient = &http.Client{
		Transport: client.transport,
		Timeout:   httpConfig.Timeout,
	}

	client.httpHandler = httpConfig.HttpHandler

	return client
}

func (client *DefaultHttpClient) SyncInvokeHttp(request *request.DefaultHttpRequest) (*response.DefaultHttpResponse,
	error) {
	exch := &exchange.SdkExchange{
		ApiReference: &exchange.ApiReference{},
		Attributes:   make(map[string]interface{}),
	}
	return client.SyncInvokeHttpWithExchange(request, exch)
}

func (client *DefaultHttpClient) SyncInvokeHttpWithExchange(request *request.DefaultHttpRequest,
	exch *exchange.SdkExchange) (*response.DefaultHttpResponse, error) {
	req, err := request.ConvertRequest()
	if err != nil {
		return nil, err
	}

	if lnErr := client.listenRequest(req); lnErr != nil {
		return nil, lnErr
	}

	client.recordRequestInfo(exch, req)
	resp, err := client.goHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	client.recordResponseInfo(exch, resp)

	if lnErr := client.listenResponse(resp); lnErr != nil {
		return nil, lnErr
	}
	client.monitorHttp(exch, resp)

	return response.NewDefaultHttpResponse(resp), nil
}

func (client *DefaultHttpClient) recordRequestInfo(exch *exchange.SdkExchange, req *http.Request) {
	exch.ApiReference.Host = req.URL.Host
	exch.ApiReference.Method = req.Method
	exch.ApiReference.Path = req.URL.Path
	exch.ApiReference.Raw = req.URL.RawQuery
	exch.ApiReference.UserAgent = req.UserAgent()
	exch.ApiReference.StartedTime = time.Now()
}

func (client *DefaultHttpClient) recordResponseInfo(exch *exchange.SdkExchange, resp *http.Response) {
	exch.ApiReference.RequestId = resp.Header.Get("X-Request-Id")
	exch.ApiReference.StatusCode = resp.StatusCode
	exch.ApiReference.ContentLength = resp.ContentLength
	exch.ApiReference.DurationMs = time.Since(exch.ApiReference.StartedTime)
}

func (client *DefaultHttpClient) listenRequest(req *http.Request) error {
	if client.httpHandler != nil && client.httpHandler.RequestHandlers != nil && req != nil {
		bodyBytes, err := httputil.DumpRequest(req, true)
		if err != nil {
			return err
		}

		reqClone := req.Clone(req.Context())
		reqClone.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		defer reqClone.Body.Close()

		client.httpHandler.RequestHandlers(*reqClone)
	}
	return nil
}

func (client *DefaultHttpClient) listenResponse(resp *http.Response) error {
	if client.httpHandler != nil && client.httpHandler.ResponseHandlers != nil && resp != nil {
		bodyBytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}

		respClone := http.Response{
			Body:             ioutil.NopCloser(bytes.NewBuffer(bodyBytes)),
			Status:           resp.Status,
			StatusCode:       resp.StatusCode,
			Proto:            resp.Proto,
			ProtoMajor:       resp.ProtoMajor,
			ProtoMinor:       resp.ProtoMinor,
			Header:           resp.Header,
			ContentLength:    resp.ContentLength,
			TransferEncoding: resp.TransferEncoding,
			Close:            resp.Close,
			Uncompressed:     resp.Uncompressed,
			Trailer:          resp.Trailer,
		}
		defer respClone.Body.Close()

		client.httpHandler.ResponseHandlers(respClone)
	}
	return nil
}

func (client *DefaultHttpClient) monitorHttp(exch *exchange.SdkExchange, resp *http.Response) {
	if client.httpHandler != nil && client.httpHandler.MonitorHandlers != nil {
		metric := &httphandler.MonitorMetric{
			Host:          exch.ApiReference.Host,
			Method:        exch.ApiReference.Method,
			Path:          exch.ApiReference.Path,
			Raw:           exch.ApiReference.Raw,
			UserAgent:     exch.ApiReference.UserAgent,
			Latency:       exch.ApiReference.DurationMs,
			RequestId:     exch.ApiReference.RequestId,
			StatusCode:    exch.ApiReference.StatusCode,
			ContentLength: exch.ApiReference.ContentLength,
		}

		client.httpHandler.MonitorHandlers(metric)
	}
}
