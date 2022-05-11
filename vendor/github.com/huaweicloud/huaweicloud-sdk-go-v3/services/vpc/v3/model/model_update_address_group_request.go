package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateAddressGroupRequest struct {

	// 地址组的唯一标识
	AddressGroupId string `json:"address_group_id"`

	Body *UpdateAddressGroupRequestBody `json:"body,omitempty"`
}

func (o UpdateAddressGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAddressGroupRequest struct{}"
	}

	return strings.Join([]string{"UpdateAddressGroupRequest", string(data)}, " ")
}
