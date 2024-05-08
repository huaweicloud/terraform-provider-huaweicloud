package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchStopInstanceRequest Request Object
type BatchStopInstanceRequest struct {
	Body *BatchShutdownInsReq `json:"body,omitempty"`
}

func (o BatchStopInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchStopInstanceRequest struct{}"
	}

	return strings.Join([]string{"BatchStopInstanceRequest", string(data)}, " ")
}
