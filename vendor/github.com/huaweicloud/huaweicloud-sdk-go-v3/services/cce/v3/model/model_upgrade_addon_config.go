package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradeAddonConfig struct {

	// 插件名称
	AddonTemplateName string `json:"addonTemplateName"`

	// 执行动作，当前升级场景支持操作为\"patch\"
	Operation string `json:"operation"`

	// 目标插件版本号
	Version string `json:"version"`

	// 插件参数列表，Key:Value格式
	Values *interface{} `json:"values,omitempty"`
}

func (o UpgradeAddonConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeAddonConfig struct{}"
	}

	return strings.Join([]string{"UpgradeAddonConfig", string(data)}, " ")
}
