package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateInstanceTopicResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateInstanceTopicResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceTopicResponse struct{}"
	}

	return strings.Join([]string{"UpdateInstanceTopicResponse", string(data)}, " ")
}
