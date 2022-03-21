package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCheckProjectPermissionForGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneCheckProjectPermissionForGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCheckProjectPermissionForGroupResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCheckProjectPermissionForGroupResponse", string(data)}, " ")
}
