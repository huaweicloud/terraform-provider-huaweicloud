package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetWtpProtectionStatusRequestInfo struct {

	// 开启关闭状态，true表示enable， false表示disable
	Status bool `json:"status"`

	// 主机ID数组，不能为空
	HostIdList []string `json:"host_id_list"`

	// 资源ID
	ResourceId *string `json:"resource_id,omitempty"`

	// 计费模式   - packet_cycle: 包周期
	ChargingMode *string `json:"charging_mode,omitempty"`
}

func (o SetWtpProtectionStatusRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetWtpProtectionStatusRequestInfo struct{}"
	}

	return strings.Join([]string{"SetWtpProtectionStatusRequestInfo", string(data)}, " ")
}
