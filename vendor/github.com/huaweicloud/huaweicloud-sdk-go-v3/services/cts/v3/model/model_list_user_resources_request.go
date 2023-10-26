package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUserResourcesRequest Request Object
type ListUserResourcesRequest struct {
}

func (o ListUserResourcesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUserResourcesRequest struct{}"
	}

	return strings.Join([]string{"ListUserResourcesRequest", string(data)}, " ")
}
