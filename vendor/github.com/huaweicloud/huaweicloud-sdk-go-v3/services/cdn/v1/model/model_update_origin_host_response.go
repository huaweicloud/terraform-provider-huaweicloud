package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateOriginHostResponse struct {
	OriginHost     *DomainOriginHost `json:"origin_host,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o UpdateOriginHostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateOriginHostResponse struct{}"
	}

	return strings.Join([]string{"UpdateOriginHostResponse", string(data)}, " ")
}
