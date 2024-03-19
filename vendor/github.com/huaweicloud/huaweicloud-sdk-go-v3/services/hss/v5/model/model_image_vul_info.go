package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ImageVulInfo struct {

	// 漏洞id
	VulId *string `json:"vul_id,omitempty"`

	// 修复紧急度，包含如下3种。   - immediate_repair ：高危。   - delay_repair ：中危。   - not_needed_repair ：低危。
	RepairNecessity *string `json:"repair_necessity,omitempty"`

	// 漏洞描述
	Description *string `json:"description,omitempty"`

	// 漏洞所在镜像层
	Position *string `json:"position,omitempty"`

	// 漏洞的软件名称
	AppName *string `json:"app_name,omitempty"`

	// 应用软件的路径（只有应用漏洞有该字段）
	AppPath *string `json:"app_path,omitempty"`

	// 软件版本
	Version *string `json:"version,omitempty"`

	// 解决方案
	Solution *string `json:"solution,omitempty"`

	// 补丁地址
	Url *string `json:"url,omitempty"`
}

func (o ImageVulInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImageVulInfo struct{}"
	}

	return strings.Join([]string{"ImageVulInfo", string(data)}, " ")
}
