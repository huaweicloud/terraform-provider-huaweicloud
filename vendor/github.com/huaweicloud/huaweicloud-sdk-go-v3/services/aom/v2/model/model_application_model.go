package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ApplicationModel struct {

	// 应用id。
	AppId *string `json:"app_id,omitempty"`

	// 应用名称。
	AppName *string `json:"app_name,omitempty"`

	// 应用来源。
	AppType *string `json:"app_type,omitempty"`
}

func (o ApplicationModel) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ApplicationModel struct{}"
	}

	return strings.Join([]string{"ApplicationModel", string(data)}, " ")
}
