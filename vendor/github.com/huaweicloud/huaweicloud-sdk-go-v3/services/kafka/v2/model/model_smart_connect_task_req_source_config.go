package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SmartConnectTaskReqSourceConfig struct {

	// Redis实例地址。（仅源端类型为Redis时需要填写）
	RedisAddress *string `json:"redis_address,omitempty"`

	// Redis实例类型。（仅源端类型为Redis时需要填写）
	RedisType *string `json:"redis_type,omitempty"`

	// DCS实例ID。（仅源端类型为Redis时需要填写）
	DcsInstanceId *string `json:"dcs_instance_id,omitempty"`

	// Redis密码。（仅源端类型为Redis时需要填写）
	RedisPassword *string `json:"redis_password,omitempty"`

	// 同步类型，“RDB_ONLY”为全量同步，“CUSTOM_OFFSET”为全量同步+增量同步。（仅源端类型为Redis时需要填写）
	SyncMode *string `json:"sync_mode,omitempty"`

	// 全量同步重试间隔时间，单位：毫秒。（仅源端类型为Redis时需要填写）
	FullSyncWaitMs *int32 `json:"full_sync_wait_ms,omitempty"`

	// 全量同步最大重试次数。（仅源端类型为Redis时需要填写）
	FullSyncMaxRetry *int32 `json:"full_sync_max_retry,omitempty"`

	// 限速，单位为KB/s。-1表示不限速。（仅源端类型为Redis时需要填写）
	Ratelimit *int32 `json:"ratelimit,omitempty"`

	// 当前Kafka实例别名。（仅源端类型为Kafka时需要填写）
	CurrentClusterName *string `json:"current_cluster_name,omitempty"`

	// 对端Kafka实例别名。（仅源端类型为Kafka时需要填写）
	ClusterName *string `json:"cluster_name,omitempty"`

	// 对端Kafka开启SASL_SSL时设置的用户名，或者创建SASL_SSL用户时设置的用户名。（仅源端类型为Kafka且对端Kafka认证方式为“SASL_SSL”时需要填写）
	UserName *string `json:"user_name,omitempty"`

	// 对端Kafka开启SASL_SSL时设置的密码，或者创建SASL_SSL用户时设置的密码。（仅源端类型为Kafka且对端Kafka认证方式为“SASL_SSL”时需要填写）
	Password *string `json:"password,omitempty"`

	// 对端Kafka认证机制。（仅源端类型为Kafka且“认证方式”为“SASL_SSL”时需要填写）
	SaslMechanism *string `json:"sasl_mechanism,omitempty"`

	// 对端Kafka实例ID。（仅源端类型为Kafka时需要填写，instance_id和bootstrap_servers仅需要填写其中一个）
	InstanceId *string `json:"instance_id,omitempty"`

	// 对端Kafka实例地址。（仅源端类型为Kafka时需要填写，instance_id和bootstrap_servers仅需要填写其中一个）
	BootstrapServers *string `json:"bootstrap_servers,omitempty"`

	// 对端Kafka认证方式。（仅源端类型为Kafka需要填写） 支持以下两种认证方式：   - SASL_SSL：表示实例已开启SASL_SSL。   - PLAINTEXT：表示实例未开启SASL_SSL。
	SecurityProtocol *string `json:"security_protocol,omitempty"`

	// 同步方向；pull为把对端Kafka实例数据复制到当前Kafka实例中，push为把当前Kafka实例数据复制到对端Kafka实例中，two-way为对两端Kafka实例数据进行双向复制。（仅源端类型为Kafka时需要填写）
	Direction *string `json:"direction,omitempty"`

	// 是否同步消费进度。（仅源端类型为Kafka时需要填写）
	SyncConsumerOffsetsEnabled *bool `json:"sync_consumer_offsets_enabled,omitempty"`

	// 在对端实例中自动创建Topic时，指定Topic的副本数，此参数值不能超过对端实例的代理数。如果对端实例中设置了“default.replication.factor”，此参数的优先级高于“default.replication.factor”。（仅源端类型为Kafka时需要填写）
	ReplicationFactor *int32 `json:"replication_factor,omitempty"`

	// 数据复制的任务数。默认值为2，建议保持默认值。如果“同步方式”为“双向”，实际任务数=设置的任务数*2。（仅源端类型为Kafka时需要填写）
	TaskNum *int32 `json:"task_num,omitempty"`

	// 是否重命名Topic，在目标Topic名称前添加源端Kafka实例的别名，形成目标Topic新的名称。（仅源端类型为Kafka时需要填写）
	RenameTopicEnabled *bool `json:"rename_topic_enabled,omitempty"`

	// 目标Topic接收复制的消息，此消息header中包含消息来源。两端实例数据双向复制时，请开启“添加来源header”，防止循环复制。（仅源端类型为Kafka时需要填写）
	ProvenanceHeaderEnabled *bool `json:"provenance_header_enabled,omitempty"`

	// 启动偏移量，latest为获取最新的数据，earliest为获取最早的数据。（仅源端类型为Kafka时需要填写）
	ConsumerStrategy *string `json:"consumer_strategy,omitempty"`

	// 复制消息所使用的压缩算法。（仅源端类型为Kafka时需要填写） - none - gzip - snappy - lz4 - zstd
	CompressionType *string `json:"compression_type,omitempty"`

	// topic映射，用于自定义目标端Topic名称。不能同时设置“重命名Topic”和“topic映射”。topic映射请按照“源端topic:目的端topic”的格式填写，如涉及多个topic映射，请用“,”分隔开，例如：topic-sc-1:topic-sc-2,topic-sc-3:topic-sc-4。（仅源端类型为Kafka时需要填写）
	TopicsMapping *string `json:"topics_mapping,omitempty"`
}

func (o SmartConnectTaskReqSourceConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmartConnectTaskReqSourceConfig struct{}"
	}

	return strings.Join([]string{"SmartConnectTaskReqSourceConfig", string(data)}, " ")
}
