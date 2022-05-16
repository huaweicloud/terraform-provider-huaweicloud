package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HlsRecordConfig struct {

	// 单位为秒，周期录制时长，最小1分钟（60秒），最大12小时。如果为0则整个流录制一个文件。
	RecordCycle int32 `json:"record_cycle"`

	// 录制m3u8文件含路径和文件名的前缀， 默认Record/{publish_domain}/{app}/{record_type}/{record_format}/{stream}_{file_start_time}/{stream}_{file_start_time}
	RecordPrefix *string `json:"record_prefix,omitempty"`

	// 录制ts文件名的前缀， 默认{file_start_time_unix}_{file_end_time_unix}_{ts_sequence_number}
	RecordTsPrefix *string `json:"record_ts_prefix,omitempty"`

	// 录制HLS时ts的切片时长，非必须，缺省为10，单位秒，最小2，最大60
	RecordSliceDuration *int32 `json:"record_slice_duration,omitempty"`

	// 录制HLS文件拼接时长，如果流中断超过该时间，则生成新文件。单位秒。如果为0表示流中断就生成新文件，如果为-1则表示相同的流中断恢复后继续在30天内的前一个文件保存。默认为0。
	RecordMaxDurationToMergeFile *int32 `json:"record_max_duration_to_merge_file,omitempty"`
}

func (o HlsRecordConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HlsRecordConfig struct{}"
	}

	return strings.Join([]string{"HlsRecordConfig", string(data)}, " ")
}
