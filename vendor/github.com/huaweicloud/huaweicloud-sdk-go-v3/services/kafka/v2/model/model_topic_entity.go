package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TopicEntity struct {

	// 是否为默认策略。
	PoliciesOnly *bool `json:"policiesOnly,omitempty"`

	// topic名称。
	Name *string `json:"name,omitempty"`

	// 副本数，配置数据的可靠性。
	Replication *int32 `json:"replication,omitempty"`

	// topic分区数，设置消费的并发数。
	Partition *int32 `json:"partition,omitempty"`

	// 消息老化时间。
	RetentionTime *int32 `json:"retention_time,omitempty"`

	// 是否开启同步复制，开启后，客户端生产消息时相应的也要设置acks=-1，否则不生效，默认关闭。
	SyncReplication *bool `json:"sync_replication,omitempty"`

	// 是否使用同步落盘。默认值为false。同步落盘会导致性能降低。
	SyncMessageFlush *bool `json:"sync_message_flush,omitempty"`

	// 扩展配置。
	ExternalConfigs *interface{} `json:"external_configs,omitempty"`

	// topic类型(0:普通Topic 1:系统(内部)Topic)。
	TopicType *int32 `json:"topic_type,omitempty"`
}

func (o TopicEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TopicEntity struct{}"
	}

	return strings.Join([]string{"TopicEntity", string(data)}, " ")
}
