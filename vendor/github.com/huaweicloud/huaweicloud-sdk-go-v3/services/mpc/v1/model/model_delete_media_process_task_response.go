package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteMediaProcessTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteMediaProcessTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMediaProcessTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteMediaProcessTaskResponse", string(data)}, " ")
}
