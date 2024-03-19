package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResetMessageOffsetWithEngineRequest Request Object
type ResetMessageOffsetWithEngineRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 消费组名称。
	Group string `json:"group"`

	Body *ResetMessageOffsetReq `json:"body,omitempty"`
}

func (o ResetMessageOffsetWithEngineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetMessageOffsetWithEngineRequest struct{}"
	}

	return strings.Join([]string{"ResetMessageOffsetWithEngineRequest", string(data)}, " ")
}
