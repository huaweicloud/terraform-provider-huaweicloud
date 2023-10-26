package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTraceResourcesResponse Response Object
type ListTraceResourcesResponse struct {

	// 返回的资源类型列表。
	Resources      *[]TraceResource `json:"resources,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListTraceResourcesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTraceResourcesResponse struct{}"
	}

	return strings.Join([]string{"ListTraceResourcesResponse", string(data)}, " ")
}
