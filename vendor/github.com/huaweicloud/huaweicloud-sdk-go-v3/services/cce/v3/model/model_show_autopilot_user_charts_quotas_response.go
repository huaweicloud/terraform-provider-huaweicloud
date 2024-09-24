package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotUserChartsQuotasResponse Response Object
type ShowAutopilotUserChartsQuotasResponse struct {
	Quotas         *QuotaRespQuotas `json:"quotas,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ShowAutopilotUserChartsQuotasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotUserChartsQuotasResponse struct{}"
	}

	return strings.Join([]string{"ShowAutopilotUserChartsQuotasResponse", string(data)}, " ")
}
