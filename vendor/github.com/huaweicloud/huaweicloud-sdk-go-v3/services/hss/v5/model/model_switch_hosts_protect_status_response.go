package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type SwitchHostsProtectStatusResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SwitchHostsProtectStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchHostsProtectStatusResponse struct{}"
	}

	return strings.Join([]string{"SwitchHostsProtectStatusResponse", string(data)}, " ")
}
