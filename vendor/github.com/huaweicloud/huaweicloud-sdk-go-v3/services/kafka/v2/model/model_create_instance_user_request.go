package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateInstanceUserRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateInstanceUserReq `json:"body,omitempty"`
}

func (o CreateInstanceUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceUserRequest struct{}"
	}

	return strings.Join([]string{"CreateInstanceUserRequest", string(data)}, " ")
}
