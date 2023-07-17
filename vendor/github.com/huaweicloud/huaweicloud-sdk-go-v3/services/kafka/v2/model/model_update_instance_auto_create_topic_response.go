package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceAutoCreateTopicResponse Response Object
type UpdateInstanceAutoCreateTopicResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateInstanceAutoCreateTopicResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceAutoCreateTopicResponse struct{}"
	}

	return strings.Join([]string{"UpdateInstanceAutoCreateTopicResponse", string(data)}, " ")
}
