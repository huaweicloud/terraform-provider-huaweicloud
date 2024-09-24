package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotQuotasResponse Response Object
type ShowAutopilotQuotasResponse struct {

	// 资源
	Quotas         *[]QuotaResource `json:"quotas,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ShowAutopilotQuotasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotQuotasResponse struct{}"
	}

	return strings.Join([]string{"ShowAutopilotQuotasResponse", string(data)}, " ")
}
