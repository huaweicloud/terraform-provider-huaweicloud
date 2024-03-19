package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateOriginHostResponse Response Object
type UpdateOriginHostResponse struct {
	OriginHost *DomainOriginHost `json:"origin_host,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateOriginHostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateOriginHostResponse struct{}"
	}

	return strings.Join([]string{"UpdateOriginHostResponse", string(data)}, " ")
}
