package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVodRetrievalRequest Request Object
type ShowVodRetrievalRequest struct {

	// 开始时间
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 采样间隔，单位：秒，取值说明： 时间跨度1天：1小时、4小时、8小时，分别对应3600秒、14400秒和28800秒。 时间跨度2~7天：1小时、4小时、8小时、1天，分别对应3600秒、14400秒、28800秒和86400秒。 时间跨度8~31天：4小时、8小时、1天，分别对应14400秒、28800秒和86400秒。 如果不传，默认取对应时间跨度的最小间隔。
	Interval *int32 `json:"interval,omitempty"`
}

func (o ShowVodRetrievalRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVodRetrievalRequest struct{}"
	}

	return strings.Join([]string{"ShowVodRetrievalRequest", string(data)}, " ")
}
