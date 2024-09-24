package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PushDomainApplication struct {

	// 应用名，默认为live
	Name *string `json:"name,omitempty"`

	// HLS切片时长，单位：s。
	HlsFragment *int32 `json:"hls_fragment,omitempty"`

	// 每个M3U8文件内ts切片个数
	HlsTsCount *int32 `json:"hls_ts_count,omitempty"`

	// 每个M3U8文件内最小ts分片数
	HlsMinFrags *int32 `json:"hls_min_frags,omitempty"`
}

func (o PushDomainApplication) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PushDomainApplication struct{}"
	}

	return strings.Join([]string{"PushDomainApplication", string(data)}, " ")
}
