package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RollingRestartResponse Response Object
type RollingRestartResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RollingRestartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RollingRestartResponse struct{}"
	}

	return strings.Join([]string{"RollingRestartResponse", string(data)}, " ")
}
