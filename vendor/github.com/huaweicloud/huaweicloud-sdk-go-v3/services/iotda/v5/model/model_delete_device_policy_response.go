package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDevicePolicyResponse Response Object
type DeleteDevicePolicyResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDevicePolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDevicePolicyResponse struct{}"
	}

	return strings.Join([]string{"DeleteDevicePolicyResponse", string(data)}, " ")
}
