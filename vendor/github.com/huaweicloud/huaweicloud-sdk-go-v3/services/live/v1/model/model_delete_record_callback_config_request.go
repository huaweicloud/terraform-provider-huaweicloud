package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteRecordCallbackConfigRequest struct {

	// 配置ID
	Id string `json:"id"`
}

func (o DeleteRecordCallbackConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRecordCallbackConfigRequest struct{}"
	}

	return strings.Join([]string{"DeleteRecordCallbackConfigRequest", string(data)}, " ")
}
