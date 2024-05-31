package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsTagsRequest Request Object
type ListClouddcnSubnetsTagsRequest struct {
}

func (o ListClouddcnSubnetsTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsTagsRequest struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsTagsRequest", string(data)}, " ")
}
