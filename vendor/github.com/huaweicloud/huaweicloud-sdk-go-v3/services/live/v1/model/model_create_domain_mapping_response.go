package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateDomainMappingResponse struct {

	// 直播播放域名
	PullDomain *string `json:"pull_domain,omitempty"`

	// 直播播放域名关联的推流域名
	PushDomain     *string `json:"push_domain,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateDomainMappingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDomainMappingResponse struct{}"
	}

	return strings.Join([]string{"CreateDomainMappingResponse", string(data)}, " ")
}
