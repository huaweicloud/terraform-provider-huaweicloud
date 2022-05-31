package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteRemuxTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteRemuxTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRemuxTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteRemuxTaskResponse", string(data)}, " ")
}
