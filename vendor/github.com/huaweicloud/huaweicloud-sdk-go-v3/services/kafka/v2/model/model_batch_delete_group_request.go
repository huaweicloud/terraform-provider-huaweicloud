package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteGroupRequest Request Object
type BatchDeleteGroupRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *BatchDeleteGroupReq `json:"body,omitempty"`
}

func (o BatchDeleteGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteGroupRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteGroupRequest", string(data)}, " ")
}
