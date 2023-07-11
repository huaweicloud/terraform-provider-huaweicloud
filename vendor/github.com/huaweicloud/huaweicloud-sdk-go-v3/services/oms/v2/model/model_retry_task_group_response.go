package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryTaskGroupResponse Response Object
type RetryTaskGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RetryTaskGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryTaskGroupResponse struct{}"
	}

	return strings.Join([]string{"RetryTaskGroupResponse", string(data)}, " ")
}
