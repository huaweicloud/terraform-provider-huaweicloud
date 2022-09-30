package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowMaintainWindowsRequest struct {
}

func (o ShowMaintainWindowsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMaintainWindowsRequest struct{}"
	}

	return strings.Join([]string{"ShowMaintainWindowsRequest", string(data)}, " ")
}
