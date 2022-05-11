package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AllowIpRangesOption struct {

	// 描述信息。
	Description *string `json:"description,omitempty"`

	// IP地址区间,例如:0.0.0.0-255.255.255.255。
	IpRange string `json:"ip_range"`
}

func (o AllowIpRangesOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowIpRangesOption struct{}"
	}

	return strings.Join([]string{"AllowIpRangesOption", string(data)}, " ")
}
