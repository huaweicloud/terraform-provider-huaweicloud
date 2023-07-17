package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteInstanceUsersRequest Request Object
type BatchDeleteInstanceUsersRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *BatchDeleteInstanceUsersReq `json:"body,omitempty"`
}

func (o BatchDeleteInstanceUsersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteInstanceUsersRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteInstanceUsersRequest", string(data)}, " ")
}
