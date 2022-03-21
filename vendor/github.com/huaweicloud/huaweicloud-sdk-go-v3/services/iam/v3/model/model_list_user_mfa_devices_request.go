package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListUserMfaDevicesRequest struct {
}

func (o ListUserMfaDevicesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUserMfaDevicesRequest struct{}"
	}

	return strings.Join([]string{"ListUserMfaDevicesRequest", string(data)}, " ")
}
