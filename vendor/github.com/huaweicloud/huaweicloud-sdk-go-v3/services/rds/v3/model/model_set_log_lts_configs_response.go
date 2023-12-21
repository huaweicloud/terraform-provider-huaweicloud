package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetLogLtsConfigsResponse Response Object
type SetLogLtsConfigsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetLogLtsConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetLogLtsConfigsResponse struct{}"
	}

	return strings.Join([]string{"SetLogLtsConfigsResponse", string(data)}, " ")
}
