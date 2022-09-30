package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type BatchDeleteInstanceUsersResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeleteInstanceUsersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteInstanceUsersResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteInstanceUsersResponse", string(data)}, " ")
}
