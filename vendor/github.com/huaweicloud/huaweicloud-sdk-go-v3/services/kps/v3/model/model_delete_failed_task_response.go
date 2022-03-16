package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteFailedTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteFailedTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteFailedTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteFailedTaskResponse", string(data)}, " ")
}
