package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDeviceProxyResponse Response Object
type DeleteDeviceProxyResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDeviceProxyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDeviceProxyResponse struct{}"
	}

	return strings.Join([]string{"DeleteDeviceProxyResponse", string(data)}, " ")
}
