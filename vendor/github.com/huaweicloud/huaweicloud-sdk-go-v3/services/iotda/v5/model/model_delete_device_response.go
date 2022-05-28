package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteDeviceResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDeviceResponse struct{}"
	}

	return strings.Join([]string{"DeleteDeviceResponse", string(data)}, " ")
}
