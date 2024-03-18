package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceUserRequest Request Object
type UpdateInstanceUserRequest struct {

	// 消息引擎的类型。
	Engine string `json:"engine"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 用户名称。
	UserName string `json:"user_name"`

	Body *UpdateUserReq `json:"body,omitempty"`
}

func (o UpdateInstanceUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceUserRequest struct{}"
	}

	return strings.Join([]string{"UpdateInstanceUserRequest", string(data)}, " ")
}
