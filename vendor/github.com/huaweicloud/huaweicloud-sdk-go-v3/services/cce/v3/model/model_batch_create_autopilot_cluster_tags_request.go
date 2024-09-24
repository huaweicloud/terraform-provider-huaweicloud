package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateAutopilotClusterTagsRequest Request Object
type BatchCreateAutopilotClusterTagsRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	Body *BatchCreateClusterTagsRequestBody `json:"body,omitempty"`
}

func (o BatchCreateAutopilotClusterTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateAutopilotClusterTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateAutopilotClusterTagsRequest", string(data)}, " ")
}
