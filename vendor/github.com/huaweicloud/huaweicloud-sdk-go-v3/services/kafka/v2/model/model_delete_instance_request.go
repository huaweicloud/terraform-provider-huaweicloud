package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteInstanceRequest Request Object
type DeleteInstanceRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o DeleteInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteInstanceRequest struct{}"
	}

	return strings.Join([]string{"DeleteInstanceRequest", string(data)}, " ")
}
