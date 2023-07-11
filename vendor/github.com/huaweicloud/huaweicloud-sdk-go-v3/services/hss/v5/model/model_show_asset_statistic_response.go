package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAssetStatisticResponse Response Object
type ShowAssetStatisticResponse struct {

	// 账号数量
	AccountNum *int64 `json:"account_num,omitempty"`

	// 开放端口数量
	PortNum *int64 `json:"port_num,omitempty"`

	// 进程数量
	ProcessNum *int64 `json:"process_num,omitempty"`

	// 软件数量
	AppNum *int64 `json:"app_num,omitempty"`

	// 自启动数量
	AutoLaunchNum *int64 `json:"auto_launch_num,omitempty"`

	// web框架数量
	WebFrameworkNum *int64 `json:"web_framework_num,omitempty"`

	// Web站点数量
	WebSiteNum *int64 `json:"web_site_num,omitempty"`

	// Jar包数量
	JarPackageNum *int64 `json:"jar_package_num,omitempty"`

	// 内核模块数量
	KernelModuleNum *int64 `json:"kernel_module_num,omitempty"`
	HttpStatusCode  int    `json:"-"`
}

func (o ShowAssetStatisticResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetStatisticResponse struct{}"
	}

	return strings.Join([]string{"ShowAssetStatisticResponse", string(data)}, " ")
}
