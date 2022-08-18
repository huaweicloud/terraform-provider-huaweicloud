package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DomainOriginHost struct {

	// 域名ID。获取方法请参见查询加速域名。
	DomainId *string `json:"domain_id,omitempty"`

	// 回源host的类型。
	OriginHostType string `json:"origin_host_type"`

	// 自定义回源host域名。
	CustomizeDomain *string `json:"customize_domain,omitempty"`
}

func (o DomainOriginHost) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainOriginHost struct{}"
	}

	return strings.Join([]string{"DomainOriginHost", string(data)}, " ")
}
