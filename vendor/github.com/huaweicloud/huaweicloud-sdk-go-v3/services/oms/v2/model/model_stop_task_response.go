package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StopTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopTaskResponse struct{}"
	}

	return strings.Join([]string{"StopTaskResponse", string(data)}, " ")
}
