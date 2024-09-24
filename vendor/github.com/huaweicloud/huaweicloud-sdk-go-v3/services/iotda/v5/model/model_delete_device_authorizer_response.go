package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDeviceAuthorizerResponse Response Object
type DeleteDeviceAuthorizerResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDeviceAuthorizerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDeviceAuthorizerResponse struct{}"
	}

	return strings.Join([]string{"DeleteDeviceAuthorizerResponse", string(data)}, " ")
}
