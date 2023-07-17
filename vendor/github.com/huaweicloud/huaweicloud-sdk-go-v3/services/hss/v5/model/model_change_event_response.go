package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeEventResponse Response Object
type ChangeEventResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeEventResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeEventResponse struct{}"
	}

	return strings.Join([]string{"ChangeEventResponse", string(data)}, " ")
}
