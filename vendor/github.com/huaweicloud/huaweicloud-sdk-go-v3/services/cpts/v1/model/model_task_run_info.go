package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TaskRunInfo TaskRunInfo
type TaskRunInfo struct {

	// 任务id
	Id *int32 `json:"id,omitempty"`

	// 任务类型（0：旧版本任务；1：新版本任务）
	RunType *int32 `json:"run_type,omitempty"`
}

func (o TaskRunInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskRunInfo struct{}"
	}

	return strings.Join([]string{"TaskRunInfo", string(data)}, " ")
}
