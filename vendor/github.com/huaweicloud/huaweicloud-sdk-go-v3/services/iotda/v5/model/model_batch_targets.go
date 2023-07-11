package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchTargets 批量操作的目标集合
type BatchTargets struct {

	// 执行批量任务的目标集合，最多支持100个目标，当task_type为firmwareUpgrade，softwareUpgrade时，此处填写device_id
	Targets *[]string `json:"targets,omitempty"`
}

func (o BatchTargets) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchTargets struct{}"
	}

	return strings.Join([]string{"BatchTargets", string(data)}, " ")
}
