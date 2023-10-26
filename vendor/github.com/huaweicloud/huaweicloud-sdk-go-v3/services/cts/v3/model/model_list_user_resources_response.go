package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUserResourcesResponse Response Object
type ListUserResourcesResponse struct {

	// 最近30天操作事件的用户列表。
	Users          *[]UserResource `json:"users,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListUserResourcesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUserResourcesResponse struct{}"
	}

	return strings.Join([]string{"ListUserResourcesResponse", string(data)}, " ")
}
