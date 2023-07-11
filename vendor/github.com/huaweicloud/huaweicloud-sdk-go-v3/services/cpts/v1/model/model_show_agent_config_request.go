package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAgentConfigRequest Request Object
type ShowAgentConfigRequest struct {
	Body *ShowAgentConfigRequestBody `json:"body,omitempty"`
}

func (o ShowAgentConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAgentConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowAgentConfigRequest", string(data)}, " ")
}
