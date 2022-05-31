package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type TagDeviceResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o TagDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagDeviceResponse struct{}"
	}

	return strings.Join([]string{"TagDeviceResponse", string(data)}, " ")
}
