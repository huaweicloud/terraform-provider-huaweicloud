package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type SetWtpProtectionStatusInfoResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetWtpProtectionStatusInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetWtpProtectionStatusInfoResponse struct{}"
	}

	return strings.Join([]string{"SetWtpProtectionStatusInfoResponse", string(data)}, " ")
}
