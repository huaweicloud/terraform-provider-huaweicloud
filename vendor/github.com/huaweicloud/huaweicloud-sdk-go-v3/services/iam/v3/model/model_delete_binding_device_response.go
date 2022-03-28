package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteBindingDeviceResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteBindingDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBindingDeviceResponse struct{}"
	}

	return strings.Join([]string{"DeleteBindingDeviceResponse", string(data)}, " ")
}
