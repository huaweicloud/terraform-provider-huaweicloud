package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateMfaDeviceReq struct {
	VirtualMfaDevice *CreateMfaDevice `json:"virtual_mfa_device"`
}

func (o CreateMfaDeviceReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMfaDeviceReq struct{}"
	}

	return strings.Join([]string{"CreateMfaDeviceReq", string(data)}, " ")
}
