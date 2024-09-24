package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotQuotasRequest Request Object
type ShowAutopilotQuotasRequest struct {
}

func (o ShowAutopilotQuotasRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotQuotasRequest struct{}"
	}

	return strings.Join([]string{"ShowAutopilotQuotasRequest", string(data)}, " ")
}
