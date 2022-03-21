package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowProjectQuotaResponse struct {
	Quotas         *QuotaResult `json:"quotas,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ShowProjectQuotaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowProjectQuotaResponse struct{}"
	}

	return strings.Join([]string{"ShowProjectQuotaResponse", string(data)}, " ")
}
