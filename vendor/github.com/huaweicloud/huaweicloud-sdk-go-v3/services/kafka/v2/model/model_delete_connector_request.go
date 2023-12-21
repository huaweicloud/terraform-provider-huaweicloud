package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConnectorRequest Request Object
type DeleteConnectorRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o DeleteConnectorRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConnectorRequest struct{}"
	}

	return strings.Join([]string{"DeleteConnectorRequest", string(data)}, " ")
}
