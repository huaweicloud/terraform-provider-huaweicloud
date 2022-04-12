package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTaskStatusRequestBody
type UpdateTaskStatusRequestBody struct {
	// cluster_id

	ClusterId int32 `json:"cluster_id"`
	// cluster_type

	ClusterType string `json:"cluster_type"`
	// without_package

	WithoutPackage int32 `json:"without_package"`

	NetworkInfo *NetworkInfo `json:"network_info"`
	// status

	Status int32 `json:"status"`
}

func (o UpdateTaskStatusRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskStatusRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTaskStatusRequestBody", string(data)}, " ")
}
