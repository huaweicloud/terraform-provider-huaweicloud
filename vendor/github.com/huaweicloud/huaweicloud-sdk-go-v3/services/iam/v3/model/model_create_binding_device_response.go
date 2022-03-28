package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateBindingDeviceResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateBindingDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateBindingDeviceResponse struct{}"
	}

	return strings.Join([]string{"CreateBindingDeviceResponse", string(data)}, " ")
}
