package auth

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"regexp"
	"strings"
)

const DefaultEndpointReg = "^[a-z][a-z0-9-]+(\\.[a-z]{2,}-[a-z]+-\\d{1,2})?\\.(my)?(huaweicloud|myhwclouds).(com|cn)"

type DerivedCredential interface {
	ProcessDerivedAuthParams(derivedAuthServiceName, regionId string) ICredential
	IsDerivedAuth(httpRequest *request.DefaultHttpRequest) bool
	ICredential
}

func GetDefaultDerivedPredicate() func(*request.DefaultHttpRequest) bool {
	return func(httpRequest *request.DefaultHttpRequest) bool {
		matched, err := regexp.MatchString(DefaultEndpointReg, strings.Replace(httpRequest.GetEndpoint(), "https://", "", 1))
		if err != nil {
			return true
		} else {
			return !matched
		}
	}
}
