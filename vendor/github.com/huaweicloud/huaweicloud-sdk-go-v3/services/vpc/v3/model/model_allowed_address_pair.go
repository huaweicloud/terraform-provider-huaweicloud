package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowedAddressPair
type AllowedAddressPair struct {

	// 功能说明：IP地址约束：不支持0.0.0.0如果allowed_address_pairs配置地址池较大的CIDR（掩码小于24位），建议为该port配置一个单独的安全组。
	IpAddress *string `json:"ip_address,omitempty"`

	// 功能说明：MAC地址
	MacAddress *string `json:"mac_address,omitempty"`
}

func (o AllowedAddressPair) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowedAddressPair struct{}"
	}

	return strings.Join([]string{"AllowedAddressPair", string(data)}, " ")
}
