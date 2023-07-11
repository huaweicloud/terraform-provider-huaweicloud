package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCoordinatorsRequest Request Object
type ShowCoordinatorsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowCoordinatorsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCoordinatorsRequest struct{}"
	}

	return strings.Join([]string{"ShowCoordinatorsRequest", string(data)}, " ")
}
