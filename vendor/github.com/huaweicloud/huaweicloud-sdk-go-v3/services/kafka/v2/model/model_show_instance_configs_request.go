package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowInstanceConfigsRequest Request Object
type ShowInstanceConfigsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowInstanceConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceConfigsRequest struct{}"
	}

	return strings.Join([]string{"ShowInstanceConfigsRequest", string(data)}, " ")
}
