package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteClouddcnSubnetsTagsResponse Response Object
type BatchDeleteClouddcnSubnetsTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeleteClouddcnSubnetsTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteClouddcnSubnetsTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteClouddcnSubnetsTagsResponse", string(data)}, " ")
}
