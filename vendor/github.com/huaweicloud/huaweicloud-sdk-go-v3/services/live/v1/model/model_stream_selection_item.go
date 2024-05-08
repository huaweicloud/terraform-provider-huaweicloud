package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StreamSelectionItem 从全量流中过滤出一个码率在[min, max]区间的流。如果不需要码率过滤可不选。
type StreamSelectionItem struct {

	// 拉流URL中用于码率过滤的参数
	Key *string `json:"key,omitempty"`

	// 最小码率，单位：bps 取值范围：0 - 104,857,600（100Mbps）
	MaxBandwidth *int32 `json:"max_bandwidth,omitempty"`

	// 最小码率，单位：bps 取值范围：0 - 104,857,600（100Mbps）
	MinBandwidth *int32 `json:"min_bandwidth,omitempty"`
}

func (o StreamSelectionItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StreamSelectionItem struct{}"
	}

	return strings.Join([]string{"StreamSelectionItem", string(data)}, " ")
}
