package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListAllProjectsPermissionsForAgencyResponse struct {

	// 权限信息列表。
	Roles *[]AgencyAllProjectRole `json:"roles,omitempty"`

	Links          *LinksSelf `json:"links,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ListAllProjectsPermissionsForAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAllProjectsPermissionsForAgencyResponse struct{}"
	}

	return strings.Join([]string{"ListAllProjectsPermissionsForAgencyResponse", string(data)}, " ")
}
