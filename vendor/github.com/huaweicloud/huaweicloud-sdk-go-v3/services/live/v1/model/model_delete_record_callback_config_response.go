package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteRecordCallbackConfigResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteRecordCallbackConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRecordCallbackConfigResponse struct{}"
	}

	return strings.Join([]string{"DeleteRecordCallbackConfigResponse", string(data)}, " ")
}
