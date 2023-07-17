package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HostVulInfoAppList struct {

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 软件版本
	AppVersion *string `json:"app_version,omitempty"`

	// 修复漏洞软件需要升级到的版本
	UpgradeVersion *string `json:"upgrade_version,omitempty"`

	// 应用软件的路径（只有应用漏洞有该字段）
	AppPath *string `json:"app_path,omitempty"`
}

func (o HostVulInfoAppList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostVulInfoAppList struct{}"
	}

	return strings.Join([]string{"HostVulInfoAppList", string(data)}, " ")
}
