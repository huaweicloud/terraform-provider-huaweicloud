package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateRecordCallbackConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateRecordCallbackConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRecordCallbackConfigResponse struct{}"
	}

	return strings.Join([]string{"CreateRecordCallbackConfigResponse", string(data)}, " ")
}
