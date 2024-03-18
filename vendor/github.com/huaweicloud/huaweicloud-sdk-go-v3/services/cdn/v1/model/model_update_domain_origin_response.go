package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainOriginResponse Response Object
type UpdateDomainOriginResponse struct {
	Origin *ResourceBody `json:"origin,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDomainOriginResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainOriginResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainOriginResponse", string(data)}, " ")
}
