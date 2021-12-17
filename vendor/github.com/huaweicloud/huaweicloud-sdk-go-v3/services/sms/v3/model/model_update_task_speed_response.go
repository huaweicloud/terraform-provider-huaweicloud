package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateTaskSpeedResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateTaskSpeedResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskSpeedResponse struct{}"
	}

	return strings.Join([]string{"UpdateTaskSpeedResponse", string(data)}, " ")
}
