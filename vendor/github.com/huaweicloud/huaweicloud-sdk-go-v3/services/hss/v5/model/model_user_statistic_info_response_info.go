package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserStatisticInfoResponseInfo 账号统计信息
type UserStatisticInfoResponseInfo struct {

	// 账号名称，参考windows文件命名规则，支持字母、数字、下划线，特殊字符!@.-等
	UserName *string `json:"user_name,omitempty"`

	// 当前账号的主机数量
	Num *int32 `json:"num,omitempty"`
}

func (o UserStatisticInfoResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserStatisticInfoResponseInfo struct{}"
	}

	return strings.Join([]string{"UserStatisticInfoResponseInfo", string(data)}, " ")
}
