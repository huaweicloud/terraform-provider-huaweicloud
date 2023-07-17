package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowInstanceUsersRequest Request Object
type ShowInstanceUsersRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowInstanceUsersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceUsersRequest struct{}"
	}

	return strings.Join([]string{"ShowInstanceUsersRequest", string(data)}, " ")
}
