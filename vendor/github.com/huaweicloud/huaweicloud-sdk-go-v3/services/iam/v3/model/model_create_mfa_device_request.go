package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateMfaDeviceRequest struct {
	Body *CreateMfaDeviceReq `json:"body,omitempty"`
}

func (o CreateMfaDeviceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMfaDeviceRequest struct{}"
	}

	return strings.Join([]string{"CreateMfaDeviceRequest", string(data)}, " ")
}
