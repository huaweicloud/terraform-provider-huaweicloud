package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CheckProjectPermissionForAgencyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CheckProjectPermissionForAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckProjectPermissionForAgencyResponse struct{}"
	}

	return strings.Join([]string{"CheckProjectPermissionForAgencyResponse", string(data)}, " ")
}
