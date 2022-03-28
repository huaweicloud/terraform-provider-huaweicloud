package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCheckroleForGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneCheckroleForGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCheckroleForGroupResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCheckroleForGroupResponse", string(data)}, " ")
}
