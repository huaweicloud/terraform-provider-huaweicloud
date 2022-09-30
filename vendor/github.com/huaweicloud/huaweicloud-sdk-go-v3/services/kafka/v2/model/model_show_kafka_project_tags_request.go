package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowKafkaProjectTagsRequest struct {
}

func (o ShowKafkaProjectTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowKafkaProjectTagsRequest struct{}"
	}

	return strings.Join([]string{"ShowKafkaProjectTagsRequest", string(data)}, " ")
}
