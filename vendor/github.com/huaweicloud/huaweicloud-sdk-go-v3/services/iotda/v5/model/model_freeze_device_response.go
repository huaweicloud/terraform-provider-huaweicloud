package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type FreezeDeviceResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o FreezeDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FreezeDeviceResponse struct{}"
	}

	return strings.Join([]string{"FreezeDeviceResponse", string(data)}, " ")
}
