package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainIp6SwitchRequest Request Object
type UpdateDomainIp6SwitchRequest struct {
	Body *DomainIpv6SwitchReq `json:"body,omitempty"`
}

func (o UpdateDomainIp6SwitchRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainIp6SwitchRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainIp6SwitchRequest", string(data)}, " ")
}
