package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteAllFailedTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteAllFailedTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAllFailedTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteAllFailedTaskResponse", string(data)}, " ")
}
