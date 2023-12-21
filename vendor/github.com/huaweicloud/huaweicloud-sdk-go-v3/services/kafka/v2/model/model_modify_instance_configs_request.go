package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyInstanceConfigsRequest Request Object
type ModifyInstanceConfigsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ModifyInstanceConfigsReq `json:"body,omitempty"`
}

func (o ModifyInstanceConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyInstanceConfigsRequest struct{}"
	}

	return strings.Join([]string{"ModifyInstanceConfigsRequest", string(data)}, " ")
}
