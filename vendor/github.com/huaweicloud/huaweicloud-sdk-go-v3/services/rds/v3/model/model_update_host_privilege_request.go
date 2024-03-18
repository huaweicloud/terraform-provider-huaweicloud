package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateHostPrivilegeRequest Request Object
type UpdateHostPrivilegeRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	Body *SetHostPrivilegeRequestV3 `json:"body,omitempty"`
}

func (o UpdateHostPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateHostPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"UpdateHostPrivilegeRequest", string(data)}, " ")
}
