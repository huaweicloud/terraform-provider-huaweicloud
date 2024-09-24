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

package internal

import (
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/impl"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/response"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
	"os"
	"reflect"
	"strings"
)

const (
	DefaultIamEndpoint         = "https://iam.myhuaweicloud.com"
	KeystoneListProjectsUri    = "/v3/projects"
	KeystoneListAuthDomainsUri = "/v3/auth/domains"
	IamEndpointEnv             = "HUAWEICLOUD_SDK_IAM_ENDPOINT"
	CreateTokenWithIdTokenUri  = "/v3.0/OS-AUTH/id-token/tokens"
	IamTraceId                 = "X-IAM-Trace-Id"
)

type IamResponse struct {
	TraceId string
}

type KeystoneListProjectsResponse struct {
	IamResponse
	Projects *[]ProjectResult `json:"projects,omitempty"`
}

type ProjectResult struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetIamEndpoint() string {
	if endpoint := os.Getenv(IamEndpointEnv); endpoint != "" {
		https := "https://"
		if !strings.HasPrefix(endpoint, https) {
			endpoint = https + endpoint
		}
		return endpoint
	}
	return DefaultIamEndpoint
}

func GetKeystoneListProjectsRequest(iamEndpoint string, regionId string, httpConfig config.HttpConfig) *request.DefaultHttpRequest {
	return request.NewHttpRequestBuilder().
		WithEndpoint(iamEndpoint).
		WithPath(KeystoneListProjectsUri).
		WithMethod("GET").
		WithSigningAlgorithm(httpConfig.SigningAlgorithm).
		AddQueryParam("name", reflect.ValueOf(regionId)).
		Build()
}

func KeystoneListProjects(client *impl.DefaultHttpClient, req *request.DefaultHttpRequest) (*KeystoneListProjectsResponse, error) {
	resp, err := client.SyncInvokeHttp(req)
	if err != nil {
		return nil, err
	}

	data, err := handleErrAndGetRespData(req, resp)
	if err != nil {
		return nil, err
	}

	keystoneListProjectResponse := new(KeystoneListProjectsResponse)
	err = utils.Unmarshal(data, keystoneListProjectResponse)
	if err != nil {
		return nil, err
	}
	keystoneListProjectResponse.TraceId = resp.GetHeader(IamTraceId)

	return keystoneListProjectResponse, nil
}

type KeystoneListAuthDomainsResponse struct {
	IamResponse
	Domains *[]Domains `json:"domains,omitempty"`
}

type Domains struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetKeystoneListAuthDomainsRequest(iamEndpoint string, httpConfig config.HttpConfig) *request.DefaultHttpRequest {
	return request.NewHttpRequestBuilder().
		WithEndpoint(iamEndpoint).
		WithPath(KeystoneListAuthDomainsUri).
		WithMethod("GET").
		WithSigningAlgorithm(httpConfig.SigningAlgorithm).
		Build()
}

func KeystoneListAuthDomains(client *impl.DefaultHttpClient, req *request.DefaultHttpRequest) (*KeystoneListAuthDomainsResponse, error) {
	resp, err := client.SyncInvokeHttp(req)
	if err != nil {
		return nil, err
	}

	data, err := handleErrAndGetRespData(req, resp)
	if err != nil {
		return nil, err
	}

	keystoneListAuthDomainsResponse := new(KeystoneListAuthDomainsResponse)
	err = utils.Unmarshal(data, keystoneListAuthDomainsResponse)
	if err != nil {
		return nil, err
	}
	keystoneListAuthDomainsResponse.TraceId = resp.GetHeader(IamTraceId)

	return keystoneListAuthDomainsResponse, nil
}

func handleErrAndGetRespData(req *request.DefaultHttpRequest, resp *response.DefaultHttpResponse) ([]byte, error) {
	if err := (sdkerr.DefaultErrorHandler{}).HandleError(req, resp); err != nil {
		return nil, err
	}

	return resp.GetBodyAsBytes()
}

type CreateTokenWithIdTokenRequest struct {
	XIdpId string                 `json:"X-Idp-Id"`
	Body   *GetIdTokenRequestBody `json:"body,omitempty"`
}

type GetIdTokenRequestBody struct {
	Auth *GetIdTokenAuthParams `json:"auth"`
}

type GetIdTokenAuthParams struct {
	IdToken *GetIdTokenIdTokenBody `json:"id_token"`

	Scope *GetIdTokenIdScopeBody `json:"scope,omitempty"`
}

type GetIdTokenIdTokenBody struct {
	Id string `json:"id"`
}

type GetIdTokenIdScopeBody struct {
	Domain *GetIdTokenScopeDomainOrProjectBody `json:"domain,omitempty"`

	Project *GetIdTokenScopeDomainOrProjectBody `json:"project,omitempty"`
}

type GetIdTokenScopeDomainOrProjectBody struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type CreateTokenWithIdTokenResponse struct {
	IamResponse
	Token          *ScopedTokenInfo `json:"token"`
	XSubjectToken  string           `json:"X-Subject-Token"`
	XRequestId     string           `json:"X-Request-Id"`
	HttpStatusCode int              `json:"-"`
}

type ScopedTokenInfo struct {
	ExpiresAt string                     `json:"expires_at"`
	Methods   []string                   `json:"methods"`
	IssuedAt  string                     `json:"issued_at"`
	User      *FederationUserBody        `json:"user"`
	Domain    *DomainInfo                `json:"domain,omitempty"`
	Project   *ProjectInfo               `json:"project,omitempty"`
	Roles     []ScopedTokenInfoRoles     `json:"roles"`
	Catalog   []UnscopedTokenInfoCatalog `json:"catalog"`
}

type FederationUserBody struct {
	OsFederation *OsFederationInfo `json:"OS-FEDERATION"`
	Domain       *DomainInfo       `json:"domain"`
	Id           *string           `json:"id,omitempty"`
	Name         *string           `json:"name,omitempty"`
}

type OsFederationInfo struct {
	IdentityProvider *IdpIdInfo      `json:"identity_provider"`
	Protocol         *ProtocolIdInfo `json:"protocol"`
	Groups           []interface{}   `json:"groups"`
}

type IdpIdInfo struct {
	Id string `json:"id"`
}

type ProtocolIdInfo struct {
	Id string `json:"id"`
}

type DomainInfo struct {
	Id   *string `json:"id,omitempty"`
	Name string  `json:"name"`
}

type ProjectInfo struct {
	Domain *DomainInfo `json:"domain,omitempty"`
	Id     *string     `json:"id,omitempty"`
	Name   string      `json:"name"`
}

type ScopedTokenInfoRoles struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type UnscopedTokenInfoCatalog struct {
	Id        *string `json:"id,omitempty"`
	Interface *string `json:"interface,omitempty"`
	Region    *string `json:"region,omitempty"`
	RegionId  *string `json:"region_id,omitempty"`
	Url       *string `json:"url,omitempty"`
}

func getCreateTokenWithIdTokenRequestBody(idToken string, scope *GetIdTokenIdScopeBody) *GetIdTokenRequestBody {
	idTokenAuth := &GetIdTokenIdTokenBody{
		Id: idToken,
	}
	authbody := &GetIdTokenAuthParams{
		IdToken: idTokenAuth,
		Scope:   scope,
	}
	body := &GetIdTokenRequestBody{
		Auth: authbody,
	}
	return body
}

func getCreateTokenWithIdTokenRequest(iamEndpoint string, idpId string, body *GetIdTokenRequestBody, httpConfig config.HttpConfig) *request.DefaultHttpRequest {
	req := request.NewHttpRequestBuilder().
		WithEndpoint(iamEndpoint).
		WithPath(CreateTokenWithIdTokenUri).
		WithMethod("POST").
		WithSigningAlgorithm(httpConfig.SigningAlgorithm).
		WithBody("body", body).
		Build()
	req.AddHeaderParam("X-Idp-Id", idpId)
	req.AddHeaderParam("Content-Type", "application/json;charset=UTF-8")
	return req
}

func GetProjectTokenWithIdTokenRequest(iamEndpoint, idpId, idToken, projectId string, httpConfig config.HttpConfig) *request.DefaultHttpRequest {
	projectScope := &GetIdTokenScopeDomainOrProjectBody{
		Id: &projectId,
	}
	scopeAuth := &GetIdTokenIdScopeBody{
		Project: projectScope,
	}
	body := getCreateTokenWithIdTokenRequestBody(idToken, scopeAuth)
	return getCreateTokenWithIdTokenRequest(iamEndpoint, idpId, body, httpConfig)
}

func GetDomainTokenWithIdTokenRequest(iamEndpoint, idpId, idToken, domainId string, httpConfig config.HttpConfig) *request.DefaultHttpRequest {
	domainScope := &GetIdTokenScopeDomainOrProjectBody{
		Id: &domainId,
	}
	scopeAuth := &GetIdTokenIdScopeBody{
		Domain: domainScope,
	}
	body := getCreateTokenWithIdTokenRequestBody(idToken, scopeAuth)
	return getCreateTokenWithIdTokenRequest(iamEndpoint, idpId, body, httpConfig)
}

func CreateTokenWithIdToken(client *impl.DefaultHttpClient, req *request.DefaultHttpRequest) (*CreateTokenWithIdTokenResponse, error) {
	resp, err := client.SyncInvokeHttp(req)
	if err != nil {
		return nil, err
	}

	data, err := handleErrAndGetRespData(req, resp)
	if err != nil {
		return nil, err
	}

	createTokenWithIdTokenResponse := new(CreateTokenWithIdTokenResponse)
	err = utils.Unmarshal(data, createTokenWithIdTokenResponse)
	if err != nil {
		return nil, err
	}

	if createTokenWithIdTokenResponse.Token.ExpiresAt == "" {
		return nil, errors.New("[CreateTokenWithIdTokenError] failed to get the expiration time of X-Auth-Token")
	}
	requestId := resp.GetHeader("X-Request-Id")
	if requestId == "" {
		return nil, errors.New("[CreateTokenWithIdTokenError] failed to get X-Request-Id")
	}
	authToken := resp.GetHeader("X-Subject-Token")
	if authToken == "" {
		return nil, errors.New("[CreateTokenWithIdTokenError] failed to get X-Auth-Token")
	}
	createTokenWithIdTokenResponse.HttpStatusCode = resp.GetStatusCode()
	createTokenWithIdTokenResponse.XRequestId = requestId
	createTokenWithIdTokenResponse.XSubjectToken = authToken
	createTokenWithIdTokenResponse.TraceId = resp.GetHeader(IamTraceId)

	return createTokenWithIdTokenResponse, nil
}
