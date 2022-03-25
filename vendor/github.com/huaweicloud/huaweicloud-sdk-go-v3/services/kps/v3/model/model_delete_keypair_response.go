package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteKeypairResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteKeypairResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteKeypairResponse struct{}"
	}

	return strings.Join([]string{"DeleteKeypairResponse", string(data)}, " ")
}
