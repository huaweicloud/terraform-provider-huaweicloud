package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TaskResult struct {

	// 媒资ID
	AssetId *string `json:"asset_id,omitempty"`

	// 修改媒资存储状态任务下发成功与否，SUCCEED成功，FAILED失败
	Status *string `json:"status,omitempty"`
}

func (o TaskResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskResult struct{}"
	}

	return strings.Join([]string{"TaskResult", string(data)}, " ")
}
