package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 切换防护的请求信息
type SwitchHostsProtectStatusRequestInfo struct {

	// 主机开通的版本，包含如下:   - hss.version.null ：无，代表关闭防护。   - hss.version.basic ：基础版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。
	Version string `json:"version"`

	// 付费模式，当version不为“hss.version.null”时，则需必填该参数   - packet_cycle : 包周期   - on_demand : 按需
	ChargingMode *string `json:"charging_mode,omitempty"`

	// HSS配额ID，不填该参数时，则随机选择对应版本配额
	ResourceId *string `json:"resource_id,omitempty"`

	// 服务器列表
	HostIdList []string `json:"host_id_list"`

	// 资源标签列表
	Tags *[]TagInfo `json:"tags,omitempty"`
}

func (o SwitchHostsProtectStatusRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchHostsProtectStatusRequestInfo struct{}"
	}

	return strings.Join([]string{"SwitchHostsProtectStatusRequestInfo", string(data)}, " ")
}
