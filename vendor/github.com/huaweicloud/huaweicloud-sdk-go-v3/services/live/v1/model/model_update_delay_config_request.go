package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDelayConfigRequest Request Object
type UpdateDelayConfigRequest struct {
	Body *ModifyDelayConfig `json:"body,omitempty"`
}

func (o UpdateDelayConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDelayConfigRequest struct{}"
	}

	return strings.Join([]string{"UpdateDelayConfigRequest", string(data)}, " ")
}
