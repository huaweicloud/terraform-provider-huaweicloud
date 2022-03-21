package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateBindingDeviceRequest struct {
	Body *BindMfaDevice `json:"body,omitempty"`
}

func (o CreateBindingDeviceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateBindingDeviceRequest struct{}"
	}

	return strings.Join([]string{"CreateBindingDeviceRequest", string(data)}, " ")
}
