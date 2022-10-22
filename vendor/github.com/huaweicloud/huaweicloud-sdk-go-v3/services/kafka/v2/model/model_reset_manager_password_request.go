package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ResetManagerPasswordRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ResetManagerPasswordReq `json:"body,omitempty"`
}

func (o ResetManagerPasswordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetManagerPasswordRequest struct{}"
	}

	return strings.Join([]string{"ResetManagerPasswordRequest", string(data)}, " ")
}
