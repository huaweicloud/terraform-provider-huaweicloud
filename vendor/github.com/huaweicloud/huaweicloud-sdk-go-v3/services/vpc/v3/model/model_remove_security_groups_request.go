package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSecurityGroupsRequest Request Object
type RemoveSecurityGroupsRequest struct {

	// 端口的唯一标识
	PortId string `json:"port_id"`

	Body *RemoveSecurityGroupsRequestBody `json:"body,omitempty"`
}

func (o RemoveSecurityGroupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSecurityGroupsRequest struct{}"
	}

	return strings.Join([]string{"RemoveSecurityGroupsRequest", string(data)}, " ")
}
