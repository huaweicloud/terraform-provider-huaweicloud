package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateLogSettingResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateLogSettingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateLogSettingResponse struct{}"
	}

	return strings.Join([]string{"UpdateLogSettingResponse", string(data)}, " ")
}
