package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTrafficMirrorSessionsResponse Response Object
type ListTrafficMirrorSessionsResponse struct {

	// 流量镜像会话
	TrafficMirrorSessions *[]TrafficMirrorSession `json:"traffic_mirror_sessions,omitempty"`

	PageInfo *PageInfo `json:"page_info,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListTrafficMirrorSessionsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrafficMirrorSessionsResponse struct{}"
	}

	return strings.Join([]string{"ListTrafficMirrorSessionsResponse", string(data)}, " ")
}
