package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CloseKafkaManagerRequest Request Object
type CloseKafkaManagerRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`
}

func (o CloseKafkaManagerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloseKafkaManagerRequest struct{}"
	}

	return strings.Join([]string{"CloseKafkaManagerRequest", string(data)}, " ")
}
