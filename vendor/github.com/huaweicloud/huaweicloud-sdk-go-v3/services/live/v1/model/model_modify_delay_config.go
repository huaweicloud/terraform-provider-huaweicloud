package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyDelayConfig struct {

	// 播放域名
	PlayDomain string `json:"play_domain"`

	// 应用名称，默认为live
	App *string `json:"app,omitempty"`

	// 延时时间，单位：ms。  包含如下取值： - 2000（低）。 - 4000（中）。 - 6000（高）。
	Delay int32 `json:"delay"`
}

func (o ModifyDelayConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyDelayConfig struct{}"
	}

	return strings.Join([]string{"ModifyDelayConfig", string(data)}, " ")
}
