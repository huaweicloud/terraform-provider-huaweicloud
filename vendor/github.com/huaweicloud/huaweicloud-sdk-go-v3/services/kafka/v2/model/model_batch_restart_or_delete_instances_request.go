package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchRestartOrDeleteInstancesRequest Request Object
type BatchRestartOrDeleteInstancesRequest struct {
	Body *BatchRestartOrDeleteInstanceReq `json:"body,omitempty"`
}

func (o BatchRestartOrDeleteInstancesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchRestartOrDeleteInstancesRequest struct{}"
	}

	return strings.Join([]string{"BatchRestartOrDeleteInstancesRequest", string(data)}, " ")
}
