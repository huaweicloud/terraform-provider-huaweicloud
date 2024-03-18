package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDatabaseOwnerRequest Request Object
type UpdateDatabaseOwnerRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *UpdateDatabaseOwnerRequestBody `json:"body,omitempty"`
}

func (o UpdateDatabaseOwnerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDatabaseOwnerRequest struct{}"
	}

	return strings.Join([]string{"UpdateDatabaseOwnerRequest", string(data)}, " ")
}
