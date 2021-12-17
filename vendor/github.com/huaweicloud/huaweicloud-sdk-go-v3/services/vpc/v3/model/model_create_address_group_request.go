package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateAddressGroupRequest struct {
	Body *CreateAddressGroupRequestBody `json:"body,omitempty"`
}

func (o CreateAddressGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAddressGroupRequest struct{}"
	}

	return strings.Join([]string{"CreateAddressGroupRequest", string(data)}, " ")
}
