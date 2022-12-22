package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type SetRaspSwitchResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetRaspSwitchResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetRaspSwitchResponse struct{}"
	}

	return strings.Join([]string{"SetRaspSwitchResponse", string(data)}, " ")
}
