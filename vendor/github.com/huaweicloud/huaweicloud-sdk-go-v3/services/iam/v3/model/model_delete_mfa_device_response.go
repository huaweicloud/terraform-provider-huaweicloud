package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteMfaDeviceResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteMfaDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMfaDeviceResponse struct{}"
	}

	return strings.Join([]string{"DeleteMfaDeviceResponse", string(data)}, " ")
}
