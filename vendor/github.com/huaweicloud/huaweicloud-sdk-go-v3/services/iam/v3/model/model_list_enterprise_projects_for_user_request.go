package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListEnterpriseProjectsForUserRequest struct {

	// 待查询用户ID。
	UserId string `json:"user_id"`
}

func (o ListEnterpriseProjectsForUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEnterpriseProjectsForUserRequest struct{}"
	}

	return strings.Join([]string{"ListEnterpriseProjectsForUserRequest", string(data)}, " ")
}
