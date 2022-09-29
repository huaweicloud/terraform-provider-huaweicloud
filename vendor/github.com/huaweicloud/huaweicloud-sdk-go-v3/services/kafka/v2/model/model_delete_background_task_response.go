package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteBackgroundTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteBackgroundTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBackgroundTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteBackgroundTaskResponse", string(data)}, " ")
}
