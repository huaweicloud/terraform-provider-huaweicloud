package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type AssociateAgencyWithProjectPermissionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AssociateAgencyWithProjectPermissionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateAgencyWithProjectPermissionResponse struct{}"
	}

	return strings.Join([]string{"AssociateAgencyWithProjectPermissionResponse", string(data)}, " ")
}
