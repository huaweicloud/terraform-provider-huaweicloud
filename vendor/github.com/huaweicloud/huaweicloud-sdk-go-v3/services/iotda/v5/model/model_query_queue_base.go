package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 队列信息结构体，返回基本信息。
type QueryQueueBase struct {

	// 队列ID，用于唯一标识一个队列。
	QueueId *string `json:"queue_id,omitempty"`

	// 队列名称，同一租户不允许重复。
	QueueName *string `json:"queue_name,omitempty"`

	// 在物联网平台创建队列的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 在物联网平台最后修改队列的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	LastModifyTime *string `json:"last_modify_time,omitempty"`
}

func (o QueryQueueBase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryQueueBase struct{}"
	}

	return strings.Join([]string{"QueryQueueBase", string(data)}, " ")
}
