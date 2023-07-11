package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowGetLogSettingResponse Response Object
type ShowGetLogSettingResponse struct {
	LogConfiguration *LogConfiguration `json:"logConfiguration,omitempty"`
	HttpStatusCode   int               `json:"-"`
}

func (o ShowGetLogSettingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGetLogSettingResponse struct{}"
	}

	return strings.Join([]string{"ShowGetLogSettingResponse", string(data)}, " ")
}
