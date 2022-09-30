package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ConfigListRsp struct {

	// 操作ID
	Id *string `json:"id,omitempty"`

	// 集群ID。
	ClusterId *string `json:"clusterId,omitempty"`

	// 创建时间。格式：Unix时间戳格式。
	CreateAt *string `json:"createAt,omitempty"`

	// 任务执行状态。 - true: 执行成功。 - false: 执行失败。
	Status *string `json:"status,omitempty"`

	// 结束时间，当创建未结束时结束时间为null。格式：Unix时间戳格式。
	FinishedAt *string `json:"finishedAt,omitempty"`

	// 修改参数配置记录。
	ModifyDeleteReset *string `json:"modifyDeleteReset,omitempty"`

	// 返回错误信息。当状态为success时该参数为null。
	FailedMsg *string `json:"failedMsg,omitempty"`
}

func (o ConfigListRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigListRsp struct{}"
	}

	return strings.Join([]string{"ConfigListRsp", string(data)}, " ")
}
