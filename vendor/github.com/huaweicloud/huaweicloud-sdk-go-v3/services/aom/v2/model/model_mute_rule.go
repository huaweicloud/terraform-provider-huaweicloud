package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 静默规则
type MuteRule struct {

	// 创建时间
	CreateTime *int64 `json:"create_time,omitempty"`

	// 规则描述
	Desc *string `json:"desc,omitempty"`

	// 规则的匹配条件
	Match [][]Match `json:"match"`

	MuteConfig *MuteConfig `json:"mute_config"`

	// 规则名称
	Name string `json:"name"`

	// 时区
	Timezone string `json:"timezone"`

	// 修改时间
	UpdateTime *int64 `json:"update_time,omitempty"`

	// 用户ID
	UserId *string `json:"user_id,omitempty"`
}

func (o MuteRule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MuteRule struct{}"
	}

	return strings.Join([]string{"MuteRule", string(data)}, " ")
}
