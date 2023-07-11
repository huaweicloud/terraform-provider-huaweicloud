package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateConnectorRequest Request Object
type CreateConnectorRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateConnectorReq `json:"body,omitempty"`
}

func (o CreateConnectorRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateConnectorRequest struct{}"
	}

	return strings.Join([]string{"CreateConnectorRequest", string(data)}, " ")
}
