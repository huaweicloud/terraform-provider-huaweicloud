package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchUpdateTasksResponse Response Object
type BatchUpdateTasksResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchUpdateTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTasksResponse struct{}"
	}

	return strings.Join([]string{"BatchUpdateTasksResponse", string(data)}, " ")
}
