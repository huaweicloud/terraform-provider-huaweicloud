package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowTagQuotaRequest struct {
}

func (o ShowTagQuotaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTagQuotaRequest struct{}"
	}

	return strings.Join([]string{"ShowTagQuotaRequest", string(data)}, " ")
}
