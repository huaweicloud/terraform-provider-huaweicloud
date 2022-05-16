package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowRecordCallbackConfigRequest struct {

	// 配置ID
	Id string `json:"id"`
}

func (o ShowRecordCallbackConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRecordCallbackConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowRecordCallbackConfigRequest", string(data)}, " ")
}
