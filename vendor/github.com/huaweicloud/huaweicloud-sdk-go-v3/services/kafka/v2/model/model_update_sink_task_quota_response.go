package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateSinkTaskQuotaResponse Response Object
type UpdateSinkTaskQuotaResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateSinkTaskQuotaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSinkTaskQuotaResponse struct{}"
	}

	return strings.Join([]string{"UpdateSinkTaskQuotaResponse", string(data)}, " ")
}
