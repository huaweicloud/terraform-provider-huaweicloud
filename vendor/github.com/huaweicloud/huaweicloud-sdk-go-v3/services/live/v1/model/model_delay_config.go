package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DelayConfig struct {

	// 应用名，默认为live
	App *string `json:"app,omitempty"`

	// 延时时间，单位：ms。  包含如下取值： - 2000（低）。 - 4000（中）。 - 6000（高）。
	Delay *int32 `json:"delay,omitempty"`
}

func (o DelayConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DelayConfig struct{}"
	}

	return strings.Join([]string{"DelayConfig", string(data)}, " ")
}
