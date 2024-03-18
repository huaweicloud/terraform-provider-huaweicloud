package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePostCheckRequest Request Object
type CreatePostCheckRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	Body *PostcheckClusterRequestBody `json:"body,omitempty"`
}

func (o CreatePostCheckRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePostCheckRequest struct{}"
	}

	return strings.Join([]string{"CreatePostCheckRequest", string(data)}, " ")
}
