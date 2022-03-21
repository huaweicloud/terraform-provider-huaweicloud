package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RemoveProjectPermissionFromAgencyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RemoveProjectPermissionFromAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveProjectPermissionFromAgencyResponse struct{}"
	}

	return strings.Join([]string{"RemoveProjectPermissionFromAgencyResponse", string(data)}, " ")
}
