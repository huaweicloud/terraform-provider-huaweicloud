package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type MfaDeviceResult struct {

	// 虚拟MFA的设备序列号。
	SerialNumber string `json:"serial_number"`

	// IAM用户ID。
	UserId string `json:"user_id"`
}

func (o MfaDeviceResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MfaDeviceResult struct{}"
	}

	return strings.Join([]string{"MfaDeviceResult", string(data)}, " ")
}
