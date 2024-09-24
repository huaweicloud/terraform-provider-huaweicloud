package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SmartConnectTaskRespSinkConfig struct {

	// Redis实例地址。（仅目标端类型为Redis时会显示）
	RedisAddress *string `json:"redis_address,omitempty"`

	// Redis实例类型。（仅目标端类型为Redis时会显示）
	RedisType *string `json:"redis_type,omitempty"`

	// DCS实例ID。（仅目标端类型为Redis时会显示）
	DcsInstanceId *string `json:"dcs_instance_id,omitempty"`

	// 目标数据库，默认为-1。（仅目标端类型为Redis时会显示）
	TargetDb *int32 `json:"target_db,omitempty"`

	// 转储启动偏移量，latest为获取最新的数据，earliest为获取最早的数据。（仅目标端类型为OBS时会显示）
	ConsumerStrategy *string `json:"consumer_strategy,omitempty"`

	// 转储文件格式。当前只支持TEXT。（仅目标端类型为OBS时会显示）
	DestinationFileType *string `json:"destination_file_type,omitempty"`

	// 记数据转储周期（秒）。（仅目标端类型为OBS时会显示）
	DeliverTimeInterval *int32 `json:"deliver_time_interval,omitempty"`

	// 转储地址。（仅目标端类型为OBS时会显示）
	ObsBucketName *string `json:"obs_bucket_name,omitempty"`

	// 转储目录。（仅目标端类型为OBS时会显示）
	ObsPath *string `json:"obs_path,omitempty"`

	// 时间目录格式。（仅目标端类型为OBS时会显示）
	PartitionFormat *string `json:"partition_format,omitempty"`

	// 记录分行符。（仅目标端类型为OBS时会显示）
	RecordDelimiter *string `json:"record_delimiter,omitempty"`

	// 存储Key。（仅目标端类型为OBS时会显示）
	StoreKeys *bool `json:"store_keys,omitempty"`

	// 每个传输文件多大后就开始上传，单位为byte；默认值5242880。（仅目标端类型为OBS时会显示）
	ObsPartSize *int32 `json:"obs_part_size,omitempty"`

	// flush_size。（仅目标端类型为OBS时会显示）
	FlushSize *int32 `json:"flush_size,omitempty"`

	// 时区。（仅目标端类型为OBS时会显示）
	Timezone *string `json:"timezone,omitempty"`

	// schema_generator类，默认为\"io.confluent.connect.storage.hive.schema.DefaultSchemaGenerator\"。（仅目标端类型为OBS时会显示）
	SchemaGeneratorClass *string `json:"schema_generator_class,omitempty"`

	// partitioner类，默认\"io.confluent.connect.storage.partitioner.TimeBasedPartitioner\"。（仅目标端类型为OBS时会显示）
	PartitionerClass *string `json:"partitioner_class,omitempty"`

	// value_converter，默认为\"org.apache.kafka.connect.converters.ByteArrayConverter\"。（仅目标端类型为OBS时会显示）
	ValueConverter *string `json:"value_converter,omitempty"`

	// key_converter，默认为\"org.apache.kafka.connect.converters.ByteArrayConverter\"。（仅目标端类型为OBS时会显示）
	KeyConverter *string `json:"key_converter,omitempty"`

	// kv_delimiter，默认为\":\"。（仅目标端类型为OBS时会显示）
	KvDelimiter *string `json:"kv_delimiter,omitempty"`
}

func (o SmartConnectTaskRespSinkConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmartConnectTaskRespSinkConfig struct{}"
	}

	return strings.Join([]string{"SmartConnectTaskRespSinkConfig", string(data)}, " ")
}
