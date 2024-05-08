package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// LiveRequestArgs 直播相关配置参数
type LiveRequestArgs struct {

	// 时延字段
	Delay *string `json:"delay,omitempty"`

	// 单位
	Unit *string `json:"unit,omitempty"`
}

func (o LiveRequestArgs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LiveRequestArgs struct{}"
	}

	return strings.Join([]string{"LiveRequestArgs", string(data)}, " ")
}
