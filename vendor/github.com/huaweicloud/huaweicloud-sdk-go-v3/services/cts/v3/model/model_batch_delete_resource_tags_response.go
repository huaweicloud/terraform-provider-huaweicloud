package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteResourceTagsResponse Response Object
type BatchDeleteResourceTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeleteResourceTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteResourceTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteResourceTagsResponse", string(data)}, " ")
}
