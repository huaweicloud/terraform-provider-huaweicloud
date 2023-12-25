package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteClusterTagsRequest Request Object
type BatchDeleteClusterTagsRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	Body *BatchDeleteClusterTagsRequestBody `json:"body,omitempty"`
}

func (o BatchDeleteClusterTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteClusterTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteClusterTagsRequest", string(data)}, " ")
}
