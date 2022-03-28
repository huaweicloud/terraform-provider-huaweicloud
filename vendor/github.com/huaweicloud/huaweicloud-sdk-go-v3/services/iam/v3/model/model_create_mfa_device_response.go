package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateMfaDeviceResponse struct {
	VirtualMfaDevice *CreateMfaDeviceRespon `json:"virtual_mfa_device,omitempty"`
	HttpStatusCode   int                    `json:"-"`
}

func (o CreateMfaDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMfaDeviceResponse struct{}"
	}

	return strings.Join([]string{"CreateMfaDeviceResponse", string(data)}, " ")
}
