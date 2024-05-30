package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateClouddcnSubnetsTagsRequest Request Object
type BatchCreateClouddcnSubnetsTagsRequest struct {

	// Clouddcn子网的id
	ResourceId string `json:"resource_id"`

	Body *BatchCreateRequestBody `json:"body,omitempty"`
}

func (o BatchCreateClouddcnSubnetsTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateClouddcnSubnetsTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateClouddcnSubnetsTagsRequest", string(data)}, " ")
}
