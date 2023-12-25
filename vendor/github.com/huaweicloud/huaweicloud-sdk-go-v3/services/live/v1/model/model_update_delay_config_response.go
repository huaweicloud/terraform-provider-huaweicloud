package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDelayConfigResponse Response Object
type UpdateDelayConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDelayConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDelayConfigResponse struct{}"
	}

	return strings.Join([]string{"UpdateDelayConfigResponse", string(data)}, " ")
}
