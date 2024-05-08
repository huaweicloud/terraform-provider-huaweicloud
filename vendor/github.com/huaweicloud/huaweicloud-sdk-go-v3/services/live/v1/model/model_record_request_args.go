package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RecordRequestArgs 录制相关配置参数
type RecordRequestArgs struct {

	// 开始时间
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 格式
	Format *string `json:"format,omitempty"`

	// 单位
	Unit *string `json:"unit,omitempty"`
}

func (o RecordRequestArgs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordRequestArgs struct{}"
	}

	return strings.Join([]string{"RecordRequestArgs", string(data)}, " ")
}
