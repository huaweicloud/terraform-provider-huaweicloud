package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetWtpProtectionStatusRequestInfo struct {

	// 开启关闭状态
	Status *bool `json:"status,omitempty"`

	// HostId list
	HostIdList *[]string `json:"host_id_list,omitempty"`

	// 资源ID
	ResourceId *string `json:"resource_id,omitempty"`

	// 随机选择配额还是指定资源
	PaymentMode *int32 `json:"payment_mode,omitempty"`
}

func (o SetWtpProtectionStatusRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetWtpProtectionStatusRequestInfo struct{}"
	}

	return strings.Join([]string{"SetWtpProtectionStatusRequestInfo", string(data)}, " ")
}
