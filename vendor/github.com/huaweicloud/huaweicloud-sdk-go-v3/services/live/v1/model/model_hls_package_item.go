package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HlsPackageItem HLS频道出流信息
type HlsPackageItem struct {

	// 客户自定义的拉流地址，包括方法、域名、路径和参数
	Url string `json:"url"`

	// 从全量流中过滤出一个码率在[min, max]区间的流。如果不需要码率过滤可不选。
	StreamSelection *[]StreamSelectionItem `json:"stream_selection,omitempty"`

	// HLS版本号
	HlsVersion *string `json:"hls_version,omitempty"`

	// 频道输出分片的时长，为必选项  单位：秒。取值范围：1-10
	SegmentDurationSeconds *int32 `json:"segment_duration_seconds,omitempty"`

	// 频道直播返回分片的窗口长度，为频道输出分片的时长乘以数量后得到的值。实际返回的分片数不小于3个。  单位：秒。取值范围：0 - 86400（24小时转化成秒后的取值）
	PlaylistWindowSeconds *int32 `json:"playlist_window_seconds,omitempty"`

	Encryption *Encryption `json:"encryption,omitempty"`

	// 广告配置
	Ads *interface{} `json:"ads,omitempty"`

	// 其他额外参数
	ExtArgs *interface{} `json:"ext_args,omitempty"`

	RequestArgs *PackageRequestArgs `json:"request_args,omitempty"`
}

func (o HlsPackageItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HlsPackageItem struct{}"
	}

	return strings.Join([]string{"HlsPackageItem", string(data)}, " ")
}
