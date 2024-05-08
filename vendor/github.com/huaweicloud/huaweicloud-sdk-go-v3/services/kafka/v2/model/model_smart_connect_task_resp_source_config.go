package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SmartConnectTaskRespSourceConfig struct {

	// Redis实例地址。（仅源端类型为Redis时会显示）
	RedisAddress *string `json:"redis_address,omitempty"`

	// Redis实例类型。（仅源端类型为Redis时会显示）
	RedisType *string `json:"redis_type,omitempty"`

	// DCS实例ID。（仅源端类型为Redis时会显示）
	DcsInstanceId *string `json:"dcs_instance_id,omitempty"`

	// 同步类型，“RDB_ONLY”为全量同步，“CUSTOM_OFFSET”为全量同步+增量同步。（仅源端类型为Redis时会显示）
	SyncMode *string `json:"sync_mode,omitempty"`

	// 全量同步重试间隔时间，单位：毫秒。（仅源端类型为Redis时会显示）
	FullSyncWaitMs *int32 `json:"full_sync_wait_ms,omitempty"`

	// 全量同步最大重试次数。（仅源端类型为Redis时会显示）
	FullSyncMaxRetry *int32 `json:"full_sync_max_retry,omitempty"`

	// 限速，单位为KB/s。-1表示不限速（仅源端类型为Redis时会显示）
	Ratelimit *int32 `json:"ratelimit,omitempty"`

	// 当前Kafka实例别名。（仅源端类型为Kafka时会显示）
	CurrentClusterName *string `json:"current_cluster_name,omitempty"`

	// 对端Kafka实例别名。（仅源端类型为Kafka时会显示）
	ClusterName *string `json:"cluster_name,omitempty"`

	// 对端Kafka用户名。（仅源端类型为Kafka时会显示）
	UserName *string `json:"user_name,omitempty"`

	// 对端Kafka认证机制。（仅源端类型为Kafka时会显示）
	SaslMechanism *string `json:"sasl_mechanism,omitempty"`

	// 对端Kafka实例ID。（仅源端类型为Kafka时会显示）
	InstanceId *string `json:"instance_id,omitempty"`

	// 对端Kafka实例地址。（仅源端类型为Kafka时会显示）
	BootstrapServers *string `json:"bootstrap_servers,omitempty"`

	// 对端Kafka认证方式。（仅源端类型为Kafka时会显示）
	SecurityProtocol *string `json:"security_protocol,omitempty"`

	// 同步方向。（仅源端类型为Kafka时会显示）
	Direction *string `json:"direction,omitempty"`

	// 是否同步消费进度。（仅源端类型为Kafka时会显示）
	SyncConsumerOffsetsEnabled *bool `json:"sync_consumer_offsets_enabled,omitempty"`

	// 副本数。（仅源端类型为Kafka时会显示）
	ReplicationFactor *int32 `json:"replication_factor,omitempty"`

	// 任务数。（仅源端类型为Kafka时会显示）
	TaskNum *int32 `json:"task_num,omitempty"`

	// 是否重命名Topic。（仅源端类型为Kafka时会显示）
	RenameTopicEnabled *bool `json:"rename_topic_enabled,omitempty"`

	// 是否添加来源header。（仅源端类型为Kafka时会显示）
	ProvenanceHeaderEnabled *bool `json:"provenance_header_enabled,omitempty"`

	// 启动偏移量，latest为获取最新的数据，earliest为获取最早的数据。（仅源端类型为Kafka时会显示）
	ConsumerStrategy *string `json:"consumer_strategy,omitempty"`

	// 压缩算法。（仅源端类型为Kafka时会显示）
	CompressionType *string `json:"compression_type,omitempty"`

	// topic映射。（仅源端类型为Kafka时会显示）
	TopicsMapping *string `json:"topics_mapping,omitempty"`
}

func (o SmartConnectTaskRespSourceConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmartConnectTaskRespSourceConfig struct{}"
	}

	return strings.Join([]string{"SmartConnectTaskRespSourceConfig", string(data)}, " ")
}
