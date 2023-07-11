package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowInstanceUsersResponse Response Object
type ShowInstanceUsersResponse struct {

	// 用户列表。
	Users          *[]ShowInstanceUsersEntity `json:"users,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ShowInstanceUsersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceUsersResponse struct{}"
	}

	return strings.Join([]string{"ShowInstanceUsersResponse", string(data)}, " ")
}
