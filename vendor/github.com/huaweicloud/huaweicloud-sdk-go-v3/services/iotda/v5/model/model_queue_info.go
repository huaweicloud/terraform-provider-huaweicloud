package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 添加时队列结构体。
type QueueInfo struct {

	// **参数说明**：队列名称，同一租户不允许重复。 **取值范围**：长度不低于8不超过128，只允许字母、数字、下划线（_）、连接符（-）、间隔号（.）、冒号（:）的组合。
	QueueName string `json:"queue_name"`
}

func (o QueueInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueueInfo struct{}"
	}

	return strings.Join([]string{"QueueInfo", string(data)}, " ")
}
