package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BlackWhiteListBody struct {

	// IP黑白名单类型（0：关闭IP黑白名单功能，1：黑名单，2：白名单）
	Type int32 `json:"type"`

	// IP黑白名单列表（支持掩码且有掩码的情况下IP必须是该IP段的第一个IP）
	IpList *[]string `json:"ip_list,omitempty"`
}

func (o BlackWhiteListBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BlackWhiteListBody struct{}"
	}

	return strings.Join([]string{"BlackWhiteListBody", string(data)}, " ")
}
