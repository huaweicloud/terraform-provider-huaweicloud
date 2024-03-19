package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IpExtraSetRespOption
type IpExtraSetRespOption struct {

	// 功能说明：单个IP地址、IP地址范围或ip地址网段，支持IPv4、IPv6
	Ip string `json:"ip"`

	// 功能说明：IP的备注信息 取值范围：0-255个字符，不能包含“<”和“>”。
	Remarks string `json:"remarks"`
}

func (o IpExtraSetRespOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IpExtraSetRespOption struct{}"
	}

	return strings.Join([]string{"IpExtraSetRespOption", string(data)}, " ")
}
