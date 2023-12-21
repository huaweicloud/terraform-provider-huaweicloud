package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowUserChartsQuotasResponse Response Object
type ShowUserChartsQuotasResponse struct {
	Quotas         *QuotaRespQuotas `json:"quotas,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ShowUserChartsQuotasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUserChartsQuotasResponse struct{}"
	}

	return strings.Join([]string{"ShowUserChartsQuotasResponse", string(data)}, " ")
}
