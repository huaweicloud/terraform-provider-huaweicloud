package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DomainOriginHost 域名回源HOST配置。
type DomainOriginHost struct {

	// 域名ID。
	DomainId *string `json:"domain_id,omitempty"`

	// 回源host的类型,accelerate：选择加速域名作为回源host域名， customize：使用自定义的域名作为回源host域名。
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
