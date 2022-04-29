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

package helper

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"os"
)

const (
	AkEnvName                     = "HUAWEICLOUD_SDK_AK"
	SkEnvName                     = "HUAWEICLOUD_SDK_SK"
	ProjectIdEnvName              = "HUAWEICLOUD_SDK_PROJECT_ID"
	DomainIdEnvName               = "HUAWEICLOUD_SDK_DOMAIN_ID"
	RegionIdEnvName               = "HUAWEICLOUD_SDK_REGION_ID"
	DerivedAuthServiceNameEnvName = "HUAWEICLOUD_SDK_DERIVED_AUTH_SERVICE_NAME"
	DerivedPredicateEnvName       = "HUAWEICLOUD_SDK_DERIVED_PREDICATE"
	DefaultDerivedPredicate       = "DEFAULT_DERIVED_PREDICATE"

	BasicCredentialType  = "basic.Credentials"
	GlobalCredentialType = "global.Credentials"
)

func LoadCredentialFromEnv(defaultType string) auth.ICredential {
	ak := os.Getenv(AkEnvName)
	sk := os.Getenv(SkEnvName)
	if ak == "" || sk == "" {
		return nil
	}

	derivedAuthServiceName := os.Getenv(DerivedAuthServiceNameEnvName)
	regionId := os.Getenv(RegionIdEnvName)

	var derivedPredicate func(*request.DefaultHttpRequest) bool
	if os.Getenv(DerivedPredicateEnvName) == DefaultDerivedPredicate {
		derivedPredicate = auth.GetDefaultDerivedPredicate()
	}

	if defaultType == BasicCredentialType {
		projectId := os.Getenv(ProjectIdEnvName)
		credential := basic.NewCredentialsBuilder().
			WithAk(ak).
			WithSk(sk).
			WithProjectId(projectId).
			WithDerivedPredicate(derivedPredicate).Build()
		return credential.ProcessDerivedAuthParams(derivedAuthServiceName, regionId)
	} else if defaultType == GlobalCredentialType {
		domainId := os.Getenv(DomainIdEnvName)
		credential := global.NewCredentialsBuilder().
			WithAk(ak).
			WithSk(sk).
			WithDomainId(domainId).
			WithDerivedPredicate(derivedPredicate).
			Build()
		return credential.ProcessDerivedAuthParams(derivedAuthServiceName, global.GlobalRegionId)
	}

	return nil
}
