package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDeleteConnectorOrderRequest Request Object
type CreateDeleteConnectorOrderRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ConnectorOrderRequestBody `json:"body,omitempty"`
}

func (o CreateDeleteConnectorOrderRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDeleteConnectorOrderRequest struct{}"
	}

	return strings.Join([]string{"CreateDeleteConnectorOrderRequest", string(data)}, " ")
}
