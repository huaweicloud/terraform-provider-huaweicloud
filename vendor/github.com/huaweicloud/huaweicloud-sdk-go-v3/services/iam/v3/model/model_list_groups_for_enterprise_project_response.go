package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListGroupsForEnterpriseProjectResponse struct {

	// 用户组信息。
	Groups         *[]ListGroupsForEnterpriseProjectResDetail `json:"groups,omitempty"`
	HttpStatusCode int                                        `json:"-"`
}

func (o ListGroupsForEnterpriseProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGroupsForEnterpriseProjectResponse struct{}"
	}

	return strings.Join([]string{"ListGroupsForEnterpriseProjectResponse", string(data)}, " ")
}
