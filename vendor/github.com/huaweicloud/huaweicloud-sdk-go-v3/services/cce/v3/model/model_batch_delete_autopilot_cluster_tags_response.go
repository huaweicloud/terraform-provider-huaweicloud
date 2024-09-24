package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteAutopilotClusterTagsResponse Response Object
type BatchDeleteAutopilotClusterTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeleteAutopilotClusterTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteAutopilotClusterTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteAutopilotClusterTagsResponse", string(data)}, " ")
}
