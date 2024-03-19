package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowStorageUsedSpaceRequest Request Object
type ShowStorageUsedSpaceRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`
}

func (o ShowStorageUsedSpaceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowStorageUsedSpaceRequest struct{}"
	}

	return strings.Join([]string{"ShowStorageUsedSpaceRequest", string(data)}, " ")
}
