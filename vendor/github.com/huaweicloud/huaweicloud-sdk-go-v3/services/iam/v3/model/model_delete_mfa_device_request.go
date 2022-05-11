package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteMfaDeviceRequest struct {

	// 绑定MFA设备的IAM 用户ID。
	UserId string `json:"user_id"`

	// MFA设备序列号。
	SerialNumber string `json:"serial_number"`
}

func (o DeleteMfaDeviceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMfaDeviceRequest struct{}"
	}

	return strings.Join([]string{"DeleteMfaDeviceRequest", string(data)}, " ")
}
