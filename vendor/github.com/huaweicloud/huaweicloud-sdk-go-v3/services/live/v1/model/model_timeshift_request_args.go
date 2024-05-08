package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TimeshiftRequestArgs 时移相关配置参数
type TimeshiftRequestArgs struct {

	// 时移时长字段名
	BackTime *string `json:"back_time,omitempty"`

	// 单位
	Unit *string `json:"unit,omitempty"`
}

func (o TimeshiftRequestArgs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TimeshiftRequestArgs struct{}"
	}

	return strings.Join([]string{"TimeshiftRequestArgs", string(data)}, " ")
}
