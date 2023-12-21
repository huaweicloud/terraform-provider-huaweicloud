package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddSecurityGroupsRequest Request Object
type AddSecurityGroupsRequest struct {

	// 端口的唯一标识
	PortId string `json:"port_id"`

	Body *AddSecurityGroupsRequestBody `json:"body,omitempty"`
}

func (o AddSecurityGroupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddSecurityGroupsRequest struct{}"
	}

	return strings.Join([]string{"AddSecurityGroupsRequest", string(data)}, " ")
}
