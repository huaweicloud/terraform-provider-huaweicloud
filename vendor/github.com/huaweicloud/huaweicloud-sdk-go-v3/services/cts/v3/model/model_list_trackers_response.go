package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTrackersResponse struct {

	// 本次查询追踪器列表返回的追踪器数组。
	Trackers       *[]TrackerResponseBody `json:"trackers,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListTrackersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrackersResponse struct{}"
	}

	return strings.Join([]string{"ListTrackersResponse", string(data)}, " ")
}
