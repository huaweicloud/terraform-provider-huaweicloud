package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UntagDeviceResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UntagDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UntagDeviceResponse struct{}"
	}

	return strings.Join([]string{"UntagDeviceResponse", string(data)}, " ")
}
