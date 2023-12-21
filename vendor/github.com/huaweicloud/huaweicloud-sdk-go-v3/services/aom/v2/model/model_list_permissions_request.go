package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPermissionsRequest Request Object
type ListPermissionsRequest struct {
}

func (o ListPermissionsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPermissionsRequest struct{}"
	}

	return strings.Join([]string{"ListPermissionsRequest", string(data)}, " ")
}
