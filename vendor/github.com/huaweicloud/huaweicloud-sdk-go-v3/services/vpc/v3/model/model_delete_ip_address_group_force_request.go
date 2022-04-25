package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteIpAddressGroupForceRequest struct {

	// IP地址组的唯一标识，要删除的IP地址组ID
	AddressGroupId string `json:"address_group_id"`
}

func (o DeleteIpAddressGroupForceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteIpAddressGroupForceRequest struct{}"
	}

	return strings.Join([]string{"DeleteIpAddressGroupForceRequest", string(data)}, " ")
}
