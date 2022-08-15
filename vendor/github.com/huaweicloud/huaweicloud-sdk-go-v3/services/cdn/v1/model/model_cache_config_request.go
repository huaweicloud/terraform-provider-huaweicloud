package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CacheConfigRequest struct {

	// 是否忽略url中的参数。
	IgnoreUrlParameter *bool `json:"ignore_url_parameter,omitempty"`

	// 缓存规则是否遵循源站
	FollowOrigin *bool `json:"follow_origin,omitempty"`

	Compress *CompressRequest `json:"compress,omitempty"`

	// 缓存规则，将覆盖之前的规则配置。规则为空重置为默认规则。
	Rules *[]Rules `json:"rules,omitempty"`
}

func (o CacheConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CacheConfigRequest struct{}"
	}

	return strings.Join([]string{"CacheConfigRequest", string(data)}, " ")
}
