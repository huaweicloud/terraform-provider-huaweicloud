package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCheckUserInGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneCheckUserInGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCheckUserInGroupResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCheckUserInGroupResponse", string(data)}, " ")
}
