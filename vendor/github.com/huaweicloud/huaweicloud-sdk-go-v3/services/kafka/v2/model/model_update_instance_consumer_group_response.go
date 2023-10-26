package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceConsumerGroupResponse Response Object
type UpdateInstanceConsumerGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateInstanceConsumerGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceConsumerGroupResponse struct{}"
	}

	return strings.Join([]string{"UpdateInstanceConsumerGroupResponse", string(data)}, " ")
}
