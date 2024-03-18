package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowOriginHostResponse Response Object
type ShowOriginHostResponse struct {
	OriginHost *DomainOriginHost `json:"origin_host,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowOriginHostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowOriginHostResponse struct{}"
	}

	return strings.Join([]string{"ShowOriginHostResponse", string(data)}, " ")
}
