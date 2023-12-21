package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeBlockedIpResponse Response Object
type ChangeBlockedIpResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeBlockedIpResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeBlockedIpResponse struct{}"
	}

	return strings.Join([]string{"ChangeBlockedIpResponse", string(data)}, " ")
}
