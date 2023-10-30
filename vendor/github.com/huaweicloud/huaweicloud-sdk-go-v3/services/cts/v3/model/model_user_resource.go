package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UserResource struct {

	// 用户名。
	Name *string `json:"name,omitempty"`
}

func (o UserResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserResource struct{}"
	}

	return strings.Join([]string{"UserResource", string(data)}, " ")
}
