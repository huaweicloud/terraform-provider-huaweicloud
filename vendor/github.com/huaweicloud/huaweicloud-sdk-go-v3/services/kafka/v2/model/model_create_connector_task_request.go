package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateConnectorTaskRequest Request Object
type CreateConnectorTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateSmartConnectTaskReq `json:"body,omitempty"`
}

func (o CreateConnectorTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateConnectorTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateConnectorTaskRequest", string(data)}, " ")
}
