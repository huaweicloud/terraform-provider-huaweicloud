package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTranscodingTaskByConsoleResponse Response Object
type DeleteTranscodingTaskByConsoleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteTranscodingTaskByConsoleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodingTaskByConsoleResponse struct{}"
	}

	return strings.Join([]string{"DeleteTranscodingTaskByConsoleResponse", string(data)}, " ")
}
