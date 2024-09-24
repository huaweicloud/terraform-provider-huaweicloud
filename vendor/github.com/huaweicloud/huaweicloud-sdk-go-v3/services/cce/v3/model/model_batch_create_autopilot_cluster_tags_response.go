package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateAutopilotClusterTagsResponse Response Object
type BatchCreateAutopilotClusterTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchCreateAutopilotClusterTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateAutopilotClusterTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateAutopilotClusterTagsResponse", string(data)}, " ")
}
