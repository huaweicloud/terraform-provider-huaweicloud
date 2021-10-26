package streams

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

type StreamDetail struct {
	AutoScaleEnabled           bool               `json:"auto_scale_enabled"`
	AutoScaleMaxPartitionCount int                `json:"auto_scale_max_partition_count"`
	AutoScaleMinPartitionCount int                `json:"auto_scale_min_partition_count"`
	CompressionFormat          string             `json:"compression_format"`
	CreateTime                 int                `json:"create_time"`
	CsvProperties              CsvProperty        `json:"csv_properties"`
	DataSchema                 string             `json:"data_schema"`
	DataType                   string             `json:"data_type"`
	LastModifiedTime           int                `json:"last_modified_time"`
	RetentionPeriod            int                `json:"retention_period"`
	Status                     string             `json:"status"` //status: CREATING,RUNNING,TERMINATING,TERMINATED
	StreamId                   string             `json:"stream_id"`
	StreamName                 string             `json:"stream_name"`
	StreamType                 string             `json:"stream_type"` //COMMON:1MB bandwidth- ADVANCED:5MB bandwidth
	Tags                       []tags.ResourceTag `json:"tags"`
	SysTags                    []tags.ResourceTag `json:"sys_tags"`
	// scaling operation record list.
	UpdatePartitionCounts []UpdatePartitionLog `json:"update_partition_counts"`
	// Total number of writable partitions (including partitions in ACTIVE state only).
	WritablePartitionCount int `json:"writable_partition_count"`
	// Total number of readable partitions (including partitions in ACTIVE and DELETED state).
	ReadablePartitionCount int `json:"readable_partition_count"`
	// A list of partitions that comprise the DIS stream.
	Partitions []Partition `json:"partitions"`
	// Specify whether there are more matching partitions of the DIS stream to list.
	HasMorePartitions bool `json:"has_more_partitions"`
}

type Partition struct {
	// Current status of each partition. CREATING,ACTIVE, DELETED, EXPIRED
	Status string `json:"status"`
	// Unique identifier of the partition.
	PartitionId string `json:"partition_id"`
	// Possible value range of the hash key used by each partition.
	HashRange string `json:"hash_range"`
	// Sequence number range of each partition.
	SequenceNumberRange string `json:"sequence_number_range"`
}

type UpdatePartitionLog struct {
	CreateTimestamp      int  `json:"create_timestamp"`
	SrcPartitionCount    int  `json:"src_partition_count"`
	TargetPartitionCount int  `json:"target_partition_count"`
	ResultCode           int  `json:"result_code"`
	ResultMsg            int  `json:"result_msg"`
	AutoScale            bool `json:"auto_scale"`
}

type ListResult struct {
	HasMoreStreams bool     `json:"has_more_streams"`
	StreamList     []Stream `json:"stream_info_list"`
	StreamNames    []string `json:"stream_names"`
	TotalNumber    int      `json:"total_number"`
}

type Stream struct {
	StreamName                 string             `json:"stream_name"`
	CreateTime                 int                `json:"create_time"`
	RetentionPeriod            int                `json:"retention_period"`
	Status                     string             `json:"status"`
	StreamType                 string             `json:"stream_type"`
	DataType                   string             `json:"data_type"`
	PartitionCount             int                `json:"partition_count"`
	AutoScaleEnabled           bool               `json:"auto_scale_enabled"`
	AutoScaleMinPartitionCount int                `json:"auto_scale_min_partition_count"`
	AutoScaleMaxPartitionCount int                `json:"auto_scale_max_partition_count"`
	Tags                       []tags.ResourceTag `json:"tags"`
	SysTags                    []tags.ResourceTag `json:"sys_tags"`
}

type ListPolicyResult struct {
	StreamId string   `json:"stream_id"`
	Rules    []Policy `json:"rules"`
}

type Policy struct {
	Principal     string `json:"principal"`
	PrincipalName string `json:"principal_name"`
	ActionType    string `json:"action_type"`
	Effect        string `json:"effect"`
}
