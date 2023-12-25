package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRegionInfoResponse Response Object
type ShowRegionInfoResponse struct {
	Body           *[]ShowRegionInfoResp `json:"body,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ShowRegionInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRegionInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowRegionInfoResponse", string(data)}, " ")
}
