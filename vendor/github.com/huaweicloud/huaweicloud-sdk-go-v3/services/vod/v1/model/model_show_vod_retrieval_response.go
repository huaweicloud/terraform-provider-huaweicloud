package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVodRetrievalResponse Response Object
type ShowVodRetrievalResponse struct {

	// 统计起始时间
	StartTime *string `json:"start_time,omitempty"`

	// 采样时间间隔
	Interval *int32 `json:"interval,omitempty"`

	SampleData     *[]VodRetrievalData `json:"sample_data,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ShowVodRetrievalResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVodRetrievalResponse struct{}"
	}

	return strings.Join([]string{"ShowVodRetrievalResponse", string(data)}, " ")
}
