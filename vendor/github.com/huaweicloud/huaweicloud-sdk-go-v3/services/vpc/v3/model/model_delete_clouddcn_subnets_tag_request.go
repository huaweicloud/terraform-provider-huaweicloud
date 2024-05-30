package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteClouddcnSubnetsTagRequest Request Object
type DeleteClouddcnSubnetsTagRequest struct {

	// Clouddcn子网的id
	ResourceId string `json:"resource_id"`

	// 待删除标签的key
	TagKey string `json:"tag_key"`
}

func (o DeleteClouddcnSubnetsTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteClouddcnSubnetsTagRequest struct{}"
	}

	return strings.Join([]string{"DeleteClouddcnSubnetsTagRequest", string(data)}, " ")
}
