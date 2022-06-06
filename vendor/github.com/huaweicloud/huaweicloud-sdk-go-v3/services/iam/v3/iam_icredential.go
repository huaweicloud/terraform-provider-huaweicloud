package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/impl"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
)

const xAuthToken = "X-Auth-Token"

type IamCredentials struct {
	AuthToken string
}

func (s *IamCredentials) ProcessAuthParams(httpClient *impl.DefaultHttpClient, region string) auth.ICredential {
	return s
}

func (s *IamCredentials) ProcessAuthRequest(httpClient *impl.DefaultHttpClient, httpRequest *request.DefaultHttpRequest) (*request.DefaultHttpRequest, error) {
	if _, ok := httpRequest.GetHeaderParams()[xAuthToken]; !ok && s.AuthToken != "" {
		httpRequest.AddHeaderParam(xAuthToken, s.AuthToken)
	}
	return httpRequest, nil
}

func NewIamCredentialsBuilder() *IamCredentialsBuilder {
	return &IamCredentialsBuilder{IamCredentials: &IamCredentials{}}
}

type IamCredentialsBuilder struct {
	IamCredentials *IamCredentials
}

func (builder *IamCredentialsBuilder) WithXAuthToken(authToken string) *IamCredentialsBuilder {
	builder.IamCredentials.AuthToken = authToken
	return builder
}

func (builder *IamCredentialsBuilder) Build() *IamCredentials {
	return builder.IamCredentials
}
