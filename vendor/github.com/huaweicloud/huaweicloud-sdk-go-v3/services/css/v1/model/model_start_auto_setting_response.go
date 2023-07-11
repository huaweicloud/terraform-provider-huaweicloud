package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartAutoSettingResponse Response Object
type StartAutoSettingResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartAutoSettingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartAutoSettingResponse struct{}"
	}

	return strings.Join([]string{"StartAutoSettingResponse", string(data)}, " ")
}
