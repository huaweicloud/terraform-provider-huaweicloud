package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsCountFilterTagsRequest Request Object
type ListClouddcnSubnetsCountFilterTagsRequest struct {
	Body *ListResourcesByTagsRequestBody `json:"body,omitempty"`
}

func (o ListClouddcnSubnetsCountFilterTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsCountFilterTagsRequest struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsCountFilterTagsRequest", string(data)}, " ")
}
