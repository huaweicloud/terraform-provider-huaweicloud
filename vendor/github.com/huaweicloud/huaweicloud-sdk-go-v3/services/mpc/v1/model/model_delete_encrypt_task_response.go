package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteEncryptTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteEncryptTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteEncryptTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteEncryptTaskResponse", string(data)}, " ")
}
