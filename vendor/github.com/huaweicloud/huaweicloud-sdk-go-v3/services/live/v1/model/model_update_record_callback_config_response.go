package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateRecordCallbackConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateRecordCallbackConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateRecordCallbackConfigResponse struct{}"
	}

	return strings.Join([]string{"UpdateRecordCallbackConfigResponse", string(data)}, " ")
}
