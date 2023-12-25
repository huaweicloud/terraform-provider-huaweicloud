package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreResult 表级时间点恢复的请求信息
type RestoreResult struct {

	// 实例ID
	InstanceId *string `json:"instance_id,omitempty"`

	// 工作流id
	JobId *string `json:"job_id,omitempty"`
}

func (o RestoreResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreResult struct{}"
	}

	return strings.Join([]string{"RestoreResult", string(data)}, " ")
}
