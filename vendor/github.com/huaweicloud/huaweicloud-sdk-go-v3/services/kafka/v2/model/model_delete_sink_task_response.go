package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSinkTaskResponse Response Object
type DeleteSinkTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteSinkTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSinkTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteSinkTaskResponse", string(data)}, " ")
}
