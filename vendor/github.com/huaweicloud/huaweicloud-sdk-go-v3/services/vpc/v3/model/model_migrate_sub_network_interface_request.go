package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type MigrateSubNetworkInterfaceRequest struct {
	Body *MigrateSubNetworkInterfaceRequestBody `json:"body,omitempty"`
}

func (o MigrateSubNetworkInterfaceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MigrateSubNetworkInterfaceRequest struct{}"
	}

	return strings.Join([]string{"MigrateSubNetworkInterfaceRequest", string(data)}, " ")
}
