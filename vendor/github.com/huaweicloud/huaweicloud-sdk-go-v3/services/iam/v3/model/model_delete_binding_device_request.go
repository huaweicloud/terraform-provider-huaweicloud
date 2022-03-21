package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteBindingDeviceRequest struct {
	Body *UnbindMfaDevice `json:"body,omitempty"`
}

func (o DeleteBindingDeviceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBindingDeviceRequest struct{}"
	}

	return strings.Join([]string{"DeleteBindingDeviceRequest", string(data)}, " ")
}
