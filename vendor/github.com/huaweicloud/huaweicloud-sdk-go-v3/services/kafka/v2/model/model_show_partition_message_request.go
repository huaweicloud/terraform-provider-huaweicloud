package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPartitionMessageRequest Request Object
type ShowPartitionMessageRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Topic名称。  Topic名称必须以字母开头且只支持大小写字母、中横线、下划线以及数字。
	Topic string `json:"topic"`

	// 分区编号。
	Partition int32 `json:"partition"`

	// 消息位置。
	MessageOffset string `json:"message_offset"`
}

func (o ShowPartitionMessageRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPartitionMessageRequest struct{}"
	}

	return strings.Join([]string{"ShowPartitionMessageRequest", string(data)}, " ")
}
