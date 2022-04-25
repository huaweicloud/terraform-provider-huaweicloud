package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowAddressGroupRequest struct {

	// 地址组的唯一标识
	AddressGroupId string `json:"address_group_id"`
}

func (o ShowAddressGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAddressGroupRequest struct{}"
	}

	return strings.Join([]string{"ShowAddressGroupRequest", string(data)}, " ")
}
