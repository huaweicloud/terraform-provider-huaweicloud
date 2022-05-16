package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type FlvRecordConfig struct {

	// 单位为秒，周期录制时长，最小1分钟，最大12小时。如果为0则整个流录制一个文件。
	RecordCycle int32 `json:"record_cycle"`

	// 录制FLV文件含路径和文件名的前缀， 默认Record/{publish_domain}/{app}/{record_type}/{record_format}/{stream}_{file_start_time}/{file_start_time}
	RecordPrefix *string `json:"record_prefix,omitempty"`

	// 录制flv拼接时长，如果流中断超过该时间，则生成新文件。单位秒。如果为0表示流中断就生成新文件。默认为0。
	RecordMaxDurationToMergeFile *int32 `json:"record_max_duration_to_merge_file,omitempty"`
}

func (o FlvRecordConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FlvRecordConfig struct{}"
	}

	return strings.Join([]string{"FlvRecordConfig", string(data)}, " ")
}
