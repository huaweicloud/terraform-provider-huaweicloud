package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRegionInfoRequest Request Object
type ShowRegionInfoRequest struct {
}

func (o ShowRegionInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRegionInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowRegionInfoRequest", string(data)}, " ")
}
