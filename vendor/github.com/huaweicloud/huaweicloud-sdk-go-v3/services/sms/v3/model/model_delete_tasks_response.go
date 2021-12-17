package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteTasksResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTasksResponse struct{}"
	}

	return strings.Join([]string{"DeleteTasksResponse", string(data)}, " ")
}
