package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceProductDataObjectInfo struct {

	// 计费模式   - packet_cycle : 包周期   - on_demand : 按需
	ChargingMode *string `json:"charging_mode,omitempty"`

	// 是否自动续费
	IsAutoRenew *bool `json:"is_auto_renew,omitempty"`

	// 版本信息,key对应的值为主机开通的版本，包含如下6种输入：   - hss.version.basic ：基础版。   - hss.version.advanced ：专业版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。   - hss.version.container.enterprise ：容器版。
	VersionInfo map[string][]ShowPeriodResponseInfo `json:"version_info,omitempty"`
}

func (o ResourceProductDataObjectInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceProductDataObjectInfo struct{}"
	}

	return strings.Join([]string{"ResourceProductDataObjectInfo", string(data)}, " ")
}
