package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResizeEngineInstanceRequest Request Object
type ResizeEngineInstanceRequest struct {

	// 消息引擎。
	Engine string `json:"engine"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ResizeEngineInstanceReq `json:"body,omitempty"`
}

func (o ResizeEngineInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeEngineInstanceRequest struct{}"
	}

	return strings.Join([]string{"ResizeEngineInstanceRequest", string(data)}, " ")
}
