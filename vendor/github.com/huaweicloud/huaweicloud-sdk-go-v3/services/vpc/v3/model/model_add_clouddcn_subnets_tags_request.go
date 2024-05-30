package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddClouddcnSubnetsTagsRequest Request Object
type AddClouddcnSubnetsTagsRequest struct {

	// Clouddcn子网的id
	ResourceId string `json:"resource_id"`

	Body *AddResourceTagsRequestBody `json:"body,omitempty"`
}

func (o AddClouddcnSubnetsTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddClouddcnSubnetsTagsRequest struct{}"
	}

	return strings.Join([]string{"AddClouddcnSubnetsTagsRequest", string(data)}, " ")
}
