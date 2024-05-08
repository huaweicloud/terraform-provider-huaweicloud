package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUsersResponse Response Object
type ListUsersResponse struct {

	// 账号总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 账号信息列表
	DataList       *[]UserResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ListUsersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUsersResponse struct{}"
	}

	return strings.Join([]string{"ListUsersResponse", string(data)}, " ")
}
