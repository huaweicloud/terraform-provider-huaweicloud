package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateClouddcnSubnetsTagsResponse Response Object
type BatchCreateClouddcnSubnetsTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchCreateClouddcnSubnetsTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateClouddcnSubnetsTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateClouddcnSubnetsTagsResponse", string(data)}, " ")
}
