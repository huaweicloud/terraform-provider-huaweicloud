package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RemoveAllProjectsPermissionFromAgencyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RemoveAllProjectsPermissionFromAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveAllProjectsPermissionFromAgencyResponse struct{}"
	}

	return strings.Join([]string{"RemoveAllProjectsPermissionFromAgencyResponse", string(data)}, " ")
}
