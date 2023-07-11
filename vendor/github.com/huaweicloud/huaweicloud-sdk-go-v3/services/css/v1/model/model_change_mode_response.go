package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeModeResponse Response Object
type ChangeModeResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeModeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeModeResponse struct{}"
	}

	return strings.Join([]string{"ChangeModeResponse", string(data)}, " ")
}
